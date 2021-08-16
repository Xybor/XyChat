package v1

import (
	"log"
	"sync"

	"github.com/xybor/xychat/models"
	xyerrors "github.com/xybor/xychat/xyerrors/v1"
)

type broadcastService struct {
	hub          chan *models.ChatMessage
	rs           roomService
	clients      []*chatService
	clientsMutex sync.Mutex
}

// Each broadcastService represents for a room in the term of broadcasting a
// message to all clients.  So there is only one broadcastService of a room in
// a time.  This map ensures that each room has only one broadcastService.
var broadcastServiceList = make(map[uint]*broadcastService)

// This mutex maintains the synchronization in reading and writing
// broadcastServiceList.
var broadcastServiceListMutex = sync.Mutex{}

// GetBroadcastService finds and returns the existed broadcastService
// associated with the given roomService.  If the broadcastService is
// non-existed, it return nil.
//
// @lock: broadcastServiceList
func GetBroadcastService(rs roomService) *broadcastService {
	broadcastServiceListMutex.Lock()
	defer broadcastServiceListMutex.Unlock()

	bService, ok := broadcastServiceList[rs.room.ID]
	if ok {
		return bService
	}

	return nil
}

// createBroadcastService creates a broadcastService from a given roomService.
// The created broadcastService will attach all current chatServices which
// associats with this room.
//
// @error: ErrorPermission
func createBroadcastService(rs roomService) (*broadcastService, xyerrors.XyError) {
	if rs.room == nil {
		return nil, xyerrors.ErrorInvalidService.New("Invalid room")
	}

	if bService := GetBroadcastService(rs); bService != nil {
		return bService, xyerrors.NoError
	}

	xerr := rs.LoadUsers()
	if xerr.Errno() != 0 {
		return nil, xerr
	}

	bService := &broadcastService{
		hub:          make(chan *models.ChatMessage),
		rs:           rs,
		clients:      make([]*chatService, 0, len(rs.room.Users)),
		clientsMutex: sync.Mutex{},
	}

	go bService.runAutoBroadcast()

	bService.findAndAttach()
	addToBroadcastList(bService)

	return bService, xyerrors.NoError
}

// addToBroadcastList adds a broadcastService to the management list.
//
// @lock: broadcastServiceList
func addToBroadcastList(bs *broadcastService) {
	broadcastServiceListMutex.Lock()
	defer broadcastServiceListMutex.Unlock()

	broadcastServiceList[bs.rs.room.ID] = bs
}

// findAndAttach finds all chatServices which joined in the room in this
// broadcastService and attach them to broadcastService.
//
// @lock: chatServiceList
func (bs *broadcastService) findAndAttach() {
	chatServiceListMutex.Lock()
	defer chatServiceListMutex.Unlock()

	for _, u := range bs.rs.room.Users {
		for _, cService := range chatServiceList[u.ID] {
			xerr := bs.attach(cService)
			if xerr.Errno() != 0 {
				log.Printf("%s, rid=%d, uid=%d\n", xerr, bs.rs.room.ID, cService.us.user.ID)
				continue
			}

			xerr = cService.connect(bs)
			if xerr.Errno() != 0 {
				log.Printf("%s, rid=%d, uid=%d\n", xerr, bs.rs.room.ID, cService.us.user.ID)

				xerr = bs.detach(cService)
				if xerr.Errno() != 0 {
					log.Panicf("%s, rid=%d, uid=%d\n", xerr, bs.rs.room.ID, cService.us.user.ID)
				}
				continue
			}
		}
	}
}

// The sleep method close the channel and remove the object from broadcastList.
// It also close the receiving channel.
// It will be permantly removed if there is no remaining reference.
//
// @lock: broadcastServiceList
func (bs *broadcastService) sleep() {
	broadcastServiceListMutex.Lock()
	defer broadcastServiceListMutex.Unlock()

	close(bs.hub)
	delete(broadcastServiceList, bs.rs.room.ID)
}

// The wakeupIfJustSleep method re-intializes the object if it had just called
// the sleep() method but there is at least one reference calling the attach()
// method.
//
// @lock: broadcastServiceList
func (bs *broadcastService) wakeupIfJustSleep() {
	broadcastServiceListMutex.Lock()
	defer broadcastServiceListMutex.Unlock()

	_, ok := broadcastServiceList[bs.rs.room.ID]
	if !ok {
		broadcastServiceList[bs.rs.room.ID] = bs
		bs.hub = make(chan *models.ChatMessage)
		go bs.runAutoBroadcast()
	}
}

// The attach method adds a chatService to the broadcastService.  If the
// broadcastService had just called the sleep() method, it will be waked up
// again.
//
// @lock: bs.clients
//
// @error: ErrorCannotAccessToRoom
func (bs *broadcastService) attach(cs *chatService) xyerrors.XyError {
	bs.clientsMutex.Lock()
	defer bs.clientsMutex.Unlock()

	if !bs.rs.contain(cs.us.user.ID) {
		return xyerrors.ErrorPermission.New(
			"User isn't entitled to be attached to this broadcast service")
	}

	alreadyAttach := false
	for _, c := range bs.clients {
		if c == cs {
			alreadyAttach = true
		}
	}

	if alreadyAttach {
		return xyerrors.NoError
	}

	bs.clients = append(bs.clients, cs)

	bs.wakeupIfJustSleep()

	return xyerrors.NoError
}

// The detach method removes a chatService from broadcastService.  If it is the
// last chatService, this method also calls the sleep() method.
//
// @lock: bs.clients
//
// @error: ErrorNotYetJoinInRoom
func (bs *broadcastService) detach(cs *chatService) xyerrors.XyError {
	bs.clientsMutex.Lock()
	defer bs.clientsMutex.Unlock()

	removedPos := -1
	for i, client := range bs.clients {
		if client == cs {
			removedPos = i
		}
	}

	if removedPos == -1 {
		return xyerrors.ErrorNotYetJoinInRoom.New(
			"User hadn't joined in room %d, cannot be deattached", bs.rs.room.ID)
	}

	// Replace the removed chatService by the last chatService.  This is the
	// optimal way to remove an element from the un-ordered array.
	bs.clients[removedPos] = bs.clients[len(bs.clients)-1]
	bs.clients = bs.clients[:len(bs.clients)-1]

	// If this is the last chatService in the broadcastService, call sleep().
	if len(bs.clients) == 0 {
		bs.sleep()
	}

	return xyerrors.NoError
}

// The broadcast method sends the received message to all clients.
//
// @lock: bs.clients
func (bs *broadcastService) broadcast(msg *models.ChatMessage) {
	bs.clientsMutex.Lock()
	defer bs.clientsMutex.Unlock()

	for _, cs := range bs.clients {
		cs.receiver <- msg
	}
}

// runAutoBroadcast waits for all messages from chatService and broadcast them.
func (bs *broadcastService) runAutoBroadcast() {
	for {
		chatMsg, ok := <-bs.hub
		if !ok {
			break
		}

		bs.broadcast(chatMsg)
	}
}
