package v1

import (
	"log"
	"sync"
	"time"

	representation "github.com/xybor/xychat/representations/v1"
)

type matchQueue struct {
	// The channel receives registered clientService
	register chan *matchService

	// The channel receives unregistered clientService
	unregister chan *matchService

	// The lock struct prevents accessing to clients map at the same time
	mutex sync.Mutex

	// The map stores all matchServices in the queue
	queue map[*matchService]bool
}

var queue *matchQueue

// InitializeMatchQueue creates a matchQueue.
//
// @goroutine: matchQueue.runRegister, matchQueue.runMatch
func InitializeMatchQueue(match_every time.Duration) {
	queue = &matchQueue{
		register:   make(chan *matchService),
		unregister: make(chan *matchService),
		queue:      map[*matchService]bool{},
	}

	go queue.runRegister()
	go queue.runMatch(match_every)
}

// GetMatchQueue returns the current matchQueue.  It is only available when the
// InitializeMatchQueue() is called.
func GetMatchQueue() *matchQueue {
	return queue
}

// GetQueue gets the matchService queue in the matchQueue struct.
//
// @note: for only debugging purpose.
func (q *matchQueue) GetQueue() map[*matchService]bool {
	return q.queue
}

// GetQueueLen returns the number of matchService in the matchQueue.
//
// @note: for only debugging purpose.
func (q *matchQueue) GetQueueLen() int {
	return len(q.queue)
}

// runRegister waits for register and unregister signals by matchServices from
// a channel and processes them.
func (q *matchQueue) runRegister() {
	for {
		select {
		case user := <-q.register:
			func() {
				q.mutex.Lock()
				defer q.mutex.Unlock()

				if _, ok := q.queue[user]; !ok {
					q.queue[user] = true
				}
			}()

		case user := <-q.unregister:
			// If a matchService unregisters, it will send zero value (invalid
			// roomid) to matchService and delete it from queue.
			func() {
				q.mutex.Lock()
				defer q.mutex.Unlock()

				if _, ok := q.queue[user]; ok {
					user.roomid <- 0
					close(user.roomid)
					delete(q.queue, user)
				}
			}()
		}
	}
}

// runMatch runs q.match after each duration.
func (q *matchQueue) runMatch(every time.Duration) {
	ticker := time.NewTicker(every)
	for {
		<-ticker.C
		q.match()
		ticker.Reset(every)
	}
}

// match tries to match all client together and sends roomid to it if a match
// is found.  If there is an error, send zero value instead.
func (q *matchQueue) match() {
	// Before matching, lock the clients map and release it after the function
	// have will finished.
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Below algorithm is a very very simple.  It simply chooses two clients in
	// turn to match until it meets the end of queue.
	var client1 *matchService = nil
	var client2 *matchService = nil

	for client := range q.queue {
		if client1 == nil {
			client1 = client
		} else if client2 == nil {
			client2 = client
		}

		if client1 != nil && client2 != nil {
			roomService := CreateRoomService(nil)

			var ID uint = 0
			if err := roomService.Create(&client1.us, &client2.us); err != nil {
				log.Println(err)
			} else {
				ID = roomService.room.ID
			}

			client1.roomid <- ID
			client2.roomid <- ID

			close(client1.roomid)
			close(client2.roomid)

			delete(q.queue, client1)
			delete(q.queue, client2)

			client1 = nil
			client2 = nil
		}
	}
}

type matchService struct {
	us userService

	// The MatchQueue in which this matchService joined
	queue *matchQueue

	// A channel receives roomid if there is a match found
	roomid chan uint

	// MatchHandler handles the room returned from the matchQueue
	MatchHandler func(representation.RoomRepresentation)
}

// A list of current existed matchServiceList with uid as the identity.
// It can't create two matchServices with the same uid.
var matchServiceList = make(map[uint]bool)

// The mutex assures that two goroutines doesn't access to matchServiceList at
// the same time.
var matchServiceListMutex = sync.Mutex{}

// CreateMatchService creates a matchService struct with a given userService.
// If there has been already a matchService with the same uid, nil will be
// returned.
//
// @lock: matchServiceList
//
// @error: ErrorPermission, ErrorDuplicatedConnection
//
// @goroutine: ms.waitForAMatch
func CreateMatchService(
	us userService,
) (*matchService, error) {
	if us.user == nil {
		return nil, ErrorPermission
	}

	matchServiceListMutex.Lock()
	defer matchServiceListMutex.Unlock()

	if _, ok := matchServiceList[us.user.ID]; ok {
		return nil, ErrorDuplicatedConnection
	}

	ms := &matchService{
		us:           us,
		queue:        GetMatchQueue(),
		roomid:       make(chan uint),
		MatchHandler: func(rr representation.RoomRepresentation) {},
	}

	matchServiceList[us.user.ID] = true

	go ms.waitForAMatch()

	return ms, nil
}

// Register push this clientService to MatchQueue
func (ms *matchService) Register() {
	ms.queue.register <- ms
}

// Unregister pop this clientService from MatchQueue
func (ms *matchService) Unregister() {
	ms.queue.unregister <- ms
}

// Close deletes the client from existed clients list.  Note that Close doesn't
// unregister from MatchQueue.
//
// @lock: matchServiceList
func (ms *matchService) Close() {
	matchServiceListMutex.Lock()
	defer matchServiceListMutex.Unlock()

	delete(matchServiceList, ms.us.user.ID)
}

// waitForAMatch waits the joinRoom signals from MatchQueue and handles the
// room with ms.MatchHandler
func (ms *matchService) waitForAMatch() {
	room := representation.RoomRepresentation{ID: <-ms.roomid}
	ms.MatchHandler(room)
}
