package v1

import "time"

type ChatMessageRepresentation struct {
	UserId    uint
	RoomId    uint
	Message   string
	CreatedAt time.Time
}
