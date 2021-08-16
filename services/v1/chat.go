package v1

import (
	"log"
	"sync"

	"github.com/xybor/xychat/models"
	r "github.com/xybor/xychat/representations/v1"
	xyerrors "github.com/xybor/xychat/xyerrors/v1"
)

// The management list of chatService.  Each user (id) can have many
// chatServices.
var chatServiceList = make(map[uint][]*chatService)
var chatServiceListMutex = sync.Mutex{}

type chatService struct {
	receiver    chan *models.ChatMessage
	us          userService
	rooms       []*broadcastService
	roomsMutex  sync.Mutex
	ChatHandler func(r.ChatMessageRepresentation) error
}

// addToChatServiceList adds a chatService to the management list.
//
// @lock: chatServiceList
func addToChatServiceList(cs *chatService) {
	chatServiceListMutex.Lock()
	defer chatServiceListMutex.Unlock()

	id := cs.us.user.ID

	if _, ok := chatServiceList[id]; !ok {
		chatServiceList[id] = make([]*chatService, 0, 3)
	}

	chatServiceList[id] = append(chatServiceList[id], cs)
}

// removeFromChatServiceList removes a chatService from the management list.
//
// @lock: chatServiceList
//
// @error: ErrorNotInManagementList
func removeFromChatServiceList(cs *chatService) xyerrors.XyError {
	chatServiceListMutex.Lock()
	defer chatServiceListMutex.Unlock()

	id := cs.us.user.ID

	if _, ok := chatServiceList[id]; !ok {
		return xyerrors.ErrorNotInManagementList.New("Cannot remove an unmanaged chat service")
	}

	for i, cService := range chatServiceList[id] {
		if cService == cs {
			shortcut := chatServiceList[id]

			if len(shortcut) == 1 {
				delete(chatServiceList, id)
			} else {
				// Remove the chatService from shortcut (a.k.a chatServiceList[id])
				shortcut[i] = shortcut[len(shortcut)-1]
				chatServiceList[id] = shortcut[:len(shortcut)-1]
			}
		}
	}

	return xyerrors.NoError
}

// CreateChatService creates a chatService from a given userService.
//
// @error: ErrorPermission
func CreateChatService(us userService) (*chatService, xyerrors.XyError) {
	us.load()

	if us.user == nil {
		return nil, xyerrors.ErrorPermission.New("User must login before")
	}

	xerr := us.LoadRooms()
	if xerr.Errno() != 0 {
		return nil, xerr
	}

	cService := &chatService{
		receiver:   make(chan *models.ChatMessage),
		us:         us,
		rooms:      make([]*broadcastService, 0, len(us.user.Rooms)),
		roomsMutex: sync.Mutex{},
	}

	return cService, xyerrors.NoError
}

// The close method closes the channel receiver, clears all rooms, and removes
// itself from chatService management list.
//
// @lock: chatServiceList
//
// @error: ErrorNotInManagementList
func (cs *chatService) close() xyerrors.XyError {
	close(cs.receiver)
	cs.rooms = nil
	return removeFromChatServiceList(cs)
}

// The wakeupIfJustSleep re-initalize the chatService.
//
// @lock: chatServiceList
/* func (cs *chatService) wakeupIfJustSleep() {
	cs.receiver = make(chan *models.ChatMessage)
	cs.rooms = make([]*broadcastService, 0, 1)
	addToChatServiceList(cs)

	go cs.runAutoReceive()
} */

// connect adds a broadcastService to this chatService.
//
// @note: it doesn't attach this chatService to the broadcastService
//
// @lock: cs.rooms, chatServiceList
//
// @error: ErrorPermission
func (cs *chatService) connect(bs *broadcastService) xyerrors.XyError {
	cs.roomsMutex.Lock()
	defer cs.roomsMutex.Unlock()

	for _, room := range cs.us.user.Rooms {
		if bs.rs.room.ID == room.ID {
			cs.rooms = append(cs.rooms, bs)
			return xyerrors.NoError
		}
	}

	return xyerrors.ErrorPermission.New("You aren't entitled to connect to room %d", bs.rs.room.ID)
}

// disconnect remove a broadcastService from the chatService.
//
// @note: It doesn't call detach() method
//
// @lock: cs.rooms, chatServiceList
//
// @error: ErrorNotYetJoinInRoom
func (cs *chatService) disconnect(bs *broadcastService, lockRooms bool) xyerrors.XyError {
	if lockRooms {
		cs.roomsMutex.Lock()
		defer cs.roomsMutex.Unlock()
	}

	pos := -1
	for i, bService := range cs.rooms {
		if bs == bService {
			pos = i
		}
	}

	if pos == -1 {
		return xyerrors.ErrorNotYetJoinInRoom.New(
			"You hadn't joined in room %d, cannot disconnect", bs.rs.room.ID)
	}

	// Remove the bs from an unordered array
	cs.rooms[pos] = cs.rooms[len(cs.rooms)-1]
	cs.rooms = cs.rooms[:len(cs.rooms)-1]

	if len(cs.rooms) == 0 {
		cs.close()
	}

	return xyerrors.NoError
}

// Online finds all broadcastServices asscociated with this chatService, then
// attachs and connects to them.
//
// @lock: bs.clients, cs.rooms, chatServiceList
func (cs *chatService) Online() {
	for _, room := range cs.us.user.Rooms {
		rs := CreateRoomService(&room.ID)

		bs, xerr := createBroadcastService(rs)

		if xerr.Errno() != 0 {
			log.Panicf("%s, rid=%d, uid=%d\n", xerr, bs.rs.room.ID, cs.us.user.ID)
		}

		xerr = bs.attach(cs)
		if xerr.Errno() != 0 {
			log.Printf("%s, rid=%d, uid=%d\n", xerr, bs.rs.room.ID, cs.us.user.ID)
			continue
		}

		xerr = cs.connect(bs)
		if xerr.Errno() != 0 {
			log.Printf("%s, rid=%d, uid=%d\n", xerr, bs.rs.room.ID, cs.us.user.ID)

			xerr = bs.detach(cs)
			if xerr.Errno() != 0 {
				log.Panicf("%s, rid=%d, uid=%d\n", xerr, bs.rs.room.ID, cs.us.user.ID)
			}

			continue
		}
	}

	addToChatServiceList(cs)
	go cs.runAutoReceive()
}

// Offline removes all broadcastServices from this chatService and calls detach
// method.
//
// @lock: cs.rooms, chatServiceList
func (cs *chatService) Offline() {
	cs.roomsMutex.Lock()
	defer cs.roomsMutex.Unlock()

	var err error

	for _, bs := range cs.rooms {
		xerr := cs.disconnect(bs, false)
		if xerr.Errno() != 0 {
			log.Panicf("%s, rid=%d, uid=%d", err, bs.rs.room.ID, cs.us.user.ID)
		}

		xerr = bs.detach(cs)
		if xerr.Errno() != 0 {
			log.Printf("%s, rid=%d, uid=%d", err, bs.rs.room.ID, cs.us.user.ID)
			continue
		}
	}
}

// The which function gets the broadcastService with the given roomid, or
// return nil if not found.
//
// @lock: cs.rooms
func (cs *chatService) which(roomid uint) *broadcastService {
	cs.roomsMutex.Lock()
	defer cs.roomsMutex.Unlock()

	for _, bs := range cs.rooms {
		if bs.rs.room.ID == roomid {
			return bs
		}
	}

	return nil
}

// SendTo send the message msg to the broadcastService with given roomid.
//
// @lock: cs.rooms
//
// @error: ErrorNotYetJoinInRoom
func (cs *chatService) SendTo(roomid uint, msg string) xyerrors.XyError {
	if bs := cs.which(roomid); bs != nil {
		chatMsg := &models.ChatMessage{RoomID: roomid, UserID: cs.us.user.ID, Message: msg}

		err := models.GetDB().Create(chatMsg).Error
		if err != nil {
			return xyerrors.ErrorUnknown
		}

		bs.hub <- chatMsg

		return xyerrors.NoError
	}

	return xyerrors.ErrorNotYetJoinInRoom.New("You hadn't joined in room %d", roomid)
}

// runAutoReceive runs the infinite loop to get messages from broadcastService
// and handle it with ChatHandler.
//
// @lock: cs.rooms, chatServiceList
func (cs *chatService) runAutoReceive() {
	for {
		msg, ok := <-cs.receiver

		if !ok {
			break
		}

		chatMsgR := r.ChatMessageRepresentation{
			UserId:    msg.UserID,
			RoomId:    msg.RoomID,
			Message:   msg.Message,
			CreatedAt: msg.CreatedAt,
		}

		err := cs.ChatHandler(chatMsgR)
		if err != nil {
			log.Println(err)
			cs.Offline()
			break
		}
	}
}
