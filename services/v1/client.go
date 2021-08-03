package v1

import (
	representation "github.com/xybor/xychat/representations/v1"
)

type clientService struct {
	user     representation.UserRepresentation
	queue    *matchQueue
	joinRoom chan uint
}

// A list of current existed clients with UID as the identity.
// It can't create two clients with the same UID.
var clients = make(map[uint]bool)

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

func (cs clientService) Register() {
	cs.queue.register <- &cs
}

func (cs clientService) Unregister() {
	cs.queue.unregister <- &cs
}

// Delete the client from existed clients list,
// then close the its channel joinRoom.
func (cs clientService) Close() {
	delete(clients, cs.user.ID)
	close(cs.joinRoom)
}

func (cs clientService) WaitForJoinRoom() representation.RoomRepresentation {
	return representation.RoomRepresentation{ID: <-cs.joinRoom}
}
