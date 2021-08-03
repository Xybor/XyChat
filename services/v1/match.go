package v1

import (
	"log"
	"sync"
	"time"
)

type matchQueue struct {
	register   chan *clientService
	unregister chan *clientService
	mutex      sync.Mutex
	clients    map[*clientService]bool
}

var queue *matchQueue

func InitializeMatchQueue() {
	queue = &matchQueue{
		register:   make(chan *clientService),
		unregister: make(chan *clientService),
		clients:    map[*clientService]bool{},
	}

	go queue.RunRegister()
	go queue.RunMatch(60 * time.Second)
}

func GetMatchQueue() *matchQueue {
	return queue
}

func (q *matchQueue) RunRegister() {
	for {
		select {
		// If the MatchQueue called Match(), the register need to stop
		// until the match process finishes.
		case user := <-q.register:
			q.mutex.Lock()
			if _, ok := q.clients[user]; !ok {
				q.clients[user] = true
			}
			q.mutex.Unlock()

		case user := <-q.unregister:
			q.mutex.Lock()
			user.joinRoom <- 0
			delete(q.clients, user)
			q.mutex.Unlock()
		}
	}
}

func (q *matchQueue) RunMatch(timeout time.Duration) {
	ticker := time.NewTicker(timeout)
	for {
		<-ticker.C
		q.Match()
		ticker.Reset(timeout)
	}
}

func (q *matchQueue) Match() {
	q.mutex.Lock()
	defer func() {
		q.mutex.Unlock()
	}()

	var client1 *clientService
	var client2 *clientService

	for client := range q.clients {
		if client1 == nil {
			client1 = client
		} else if client2 == nil {
			client2 = client
		}

		if client1 != nil && client2 != nil {
			rs := roomService{}

			var ID uint = 0
			if err := rs.Create(); err != nil {
				log.Println(err)
			} else {
				ID = *rs.id
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
