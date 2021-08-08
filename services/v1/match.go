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

// InitializeMatchQueue creates a matchQueue and runs two goroutines.  The
// matchQueue.runRegister will receive all register and unregister signals.
// The matchQueue.runMatch will run matching algorithms every N seconds.
func InitializeMatchQueue(match_every time.Duration) {
	queue = &matchQueue{
		register:   make(chan *matchService),
		unregister: make(chan *matchService),
		queue:      map[*matchService]bool{},
	}

	go queue.runRegister()
	go queue.runMatch(match_every)
}

// GetMatchQueue returns the current matchQueue
func GetMatchQueue() *matchQueue {
	return queue
}

func (q *matchQueue) GetQueue() map[*matchService]bool {
	return q.queue
}

func (q *matchQueue) GetQueueLen() int {
	return len(q.queue)
}

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
					user.joinRoom <- 0
					close(user.joinRoom)
					delete(q.queue, user)
				}
			}()
		}
	}
}

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
	var client1 *matchService
	var client2 *matchService

	for client := range q.queue {
		if client1 == nil {
			client1 = client
		} else if client2 == nil {
			client2 = client
		}

		if client1 != nil && client2 != nil {
			rservice := CreateRoomService(nil)

			var ID uint = 0
			if err := rservice.create(); err != nil {
				log.Println(err)
			} else {
				ID = rservice.room.ID
			}

			client1.joinRoom <- ID
			client2.joinRoom <- ID

			close(client1.joinRoom)
			close(client2.joinRoom)

			delete(q.queue, client1)
			delete(q.queue, client2)

			client1 = nil
			client2 = nil
		}
	}
}

type matchService struct {
	userService

	// The MatchQueue will be sent to a
	queue *matchQueue

	// A channel receives roomid if there is a match
	joinRoom chan uint
	room     chan representation.RoomRepresentation
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
		userService: us,
		queue:       GetMatchQueue(),
		joinRoom:    make(chan uint),
		room:        make(chan representation.RoomRepresentation),
	}

	matchServiceList[us.user.ID] = true

	go ms.waitForJoinRoom()

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
func (ms *matchService) Close() {
	matchServiceListMutex.Lock()
	defer matchServiceListMutex.Unlock()
	delete(matchServiceList, ms.user.ID)
}

func (ms *matchService) waitForJoinRoom() {
	ms.room <- representation.RoomRepresentation{ID: <-ms.joinRoom}
	close(ms.room)
}

// WaitForRoom waits the value from a channel.  If MatchQueue finds a
// match, it will send the roomid via this channel. If an error occurs, a zero
// value will be sent. The returned value type is RoomRepresentation.
func (ms *matchService) WaitForRoom() representation.RoomRepresentation {
	return <-ms.room
}
