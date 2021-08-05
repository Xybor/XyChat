package v1

import (
	representation "github.com/xybor/xychat/representations/v1"
)

type clientService struct {
	user representation.UserRepresentation

	// The MatchQueue will be sent to a
	queue *matchQueue

	// A channel receives roomid if there is a match
	joinRoom chan uint
}

// A list of current existed clients with uid as the identity.
// It can't create two client services with the same uid.
var clients = make(map[uint]bool)

// CreateClientService creates a clientService struct with a given
// UserRepresentation.  If there has been already a clientService with the
// same uid, nil will be returned.
func CreateClientService(
	user representation.UserRepresentation,
) *clientService {
	if _, ok := clients[user.ID]; !ok {
		return nil
	}

	cs := &clientService{
		user:     user,
		queue:    GetMatchQueue(),
		joinRoom: make(chan uint),
	}

	clients[user.ID] = true

	return cs
}

// Register push this clientService to MatchQueue
func (cs *clientService) Register() {
	cs.queue.register <- cs
}

// Unregister pop this clientService from MatchQueue
func (cs *clientService) Unregister() {
	cs.queue.unregister <- cs
}

// Close deletes the client from existed clients list, then close the its
// channel joinRoom.  Note that Close doesn't unregister from MatchQueue.
func (cs *clientService) Close() {
	delete(clients, cs.user.ID)
	close(cs.joinRoom)
}

// WaitForJoinRoom waits the value from a channel.  If MatchQueue finds a
// match, it will send the roomid via this channel. If an error occurs, a zero
// value will be sent. The returned value type is RoomRepresentation.
func (cs *clientService) WaitForJoinRoom() representation.RoomRepresentation {
	return representation.RoomRepresentation{ID: <-cs.joinRoom}
}
