package v1

import "time"

type ChatMessageResponse struct {
	UserId  uint      `json:"userid"`
	RoomId  uint      `json:"roomid"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type ChatMessageRequest struct {
	RoomId  uint   `json:"roomid"`
	Message string `json:"message"`
}
