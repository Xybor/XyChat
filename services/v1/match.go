package v1

import (
	"log"
	"sync"
	"time"
)

type matchQueue struct {
	// The channel receives registered clientService
	register chan *clientService

	// The channel receives unregistered clientService
	unregister chan *clientService

	// The lock struct prevents accessing to clients map at the same time
	mutex sync.Mutex

	// The map stores all clientServices in the queue
	clients map[*clientService]bool
}

var queue *matchQueue

// InitializeMatchQueue creates a MatchQueue and runs two goroutines.  The
// matchQueue.runRegister will receive all register and unregister signal.  The
// matchQueue.runMatch will run matching algorithms every N seconds.
func InitializeMatchQueue() {
	queue = &matchQueue{
		register:   make(chan *clientService),
		unregister: make(chan *clientService),
		clients:    map[*clientService]bool{},
	}

	go queue.runRegister()
	go queue.runMatch(60 * time.Second)
}

// GetMatchQueue returns the current matchQueue
func GetMatchQueue() *matchQueue {
	return queue
}

func (q *matchQueue) runRegister() {
	for {
		select {
		case user := <-q.register:
			q.mutex.Lock()
			if _, ok := q.clients[user]; !ok {
				q.clients[user] = true
			}
			q.mutex.Unlock()

		case user := <-q.unregister:
			// If a clientService unregisters, it will send zero value (invalid
			// roomid) to clientService and delete it from queue.
			q.mutex.Lock()
			user.joinRoom <- 0
			delete(q.clients, user)
			q.mutex.Unlock()
		}
	}
}

func (q *matchQueue) runMatch(timeout time.Duration) {
	ticker := time.NewTicker(timeout)
	for {
		<-ticker.C
		q.match()
		ticker.Reset(timeout)
	}
}

// match tries to match all client together and sends roomid to it if a match
// is found.  If there is an error, send zero value instead.
func (q *matchQueue) match() {
	// Before matching, lock the clients map and release it after the function
	// have will finished.
	q.mutex.Lock()
	defer func() {
		q.mutex.Unlock()
	}()


	// Below algorithm is a very very simple.  It simply chooses two clients
	// in turn to match until it meets the end of queue.
	var client1 *clientService
	var client2 *clientService

	for client := range q.clients {
		if client1 == nil {
			client1 = client
		} else if client2 == nil {
			client2 = client
		}

		if client1 != nil && client2 != nil {
			rservice := roomService{}

			var ID uint = 0
			if err := rservice.Create(); err != nil {
				log.Println(err)
			} else {
				ID = *rservice.id
			}

			client1.joinRoom <- ID
			client2.joinRoom <- ID

			delete(q.clients, client1)
			delete(q.clients, client2)

			client1 = nil
			client2 = nil
		}
	}
}
