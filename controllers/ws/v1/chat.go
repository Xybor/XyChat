package v1

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xybor/xychat/helpers/context"
	wshelpers "github.com/xybor/xychat/helpers/ws/v1"
	r "github.com/xybor/xychat/representations/v1"
	services "github.com/xybor/xychat/services/v1"
	xyerrors "github.com/xybor/xychat/xyerrors/v1"
)

type clientMessage struct {
	RoomId  uint   `json:"roomid"`
	Message string `json:"message"`
}

func WSChatHandler(ctx *gin.Context) {
	connection := ctx.MustGet("WebSocket")
	conn := connection.(*websocket.Conn)
	client := CreateWSClient(conn)

	uid := context.GetUID(ctx)
	userService := services.CreateUserService(uid, true)

	chatService, xerr := services.CreateChatService(userService)
	if xerr.Errno() != 0 {
		response := wshelpers.NewWSError(xerr)
		client.WriteJSON(response)
		return
	}

	var isAlive = make(chan bool)

	client.CloseHandler = func() {
		chatService.Offline()

		client.Close()

		isAlive <- false
		close(isAlive)
	}

	// Receive a message from the client and send it to the corresponding room.
	client.ReadHandler = func(message []byte) error {
		msg := clientMessage{}
		err := json.Unmarshal(message, &msg)
		if err != nil {
			response := wshelpers.NewWSError(xyerrors.ErrorSyntaxInput.New("invalid json format"))
			client.WriteJSON(response)
			return nil
		}

		xerr = chatService.SendTo(msg.RoomId, msg.Message)
		if xerr.Errno() != 0 {
			response := wshelpers.NewWSError(xerr)
			client.WriteJSON(response)
			return nil
		}

		return nil
	}

	// Forward the message from broadcast service to client
	chatService.ChatHandler = func(cmr r.ChatMessageRepresentation) error {
		response := wshelpers.NewWSResponse(cmr)
		client.WriteJSON(response)
		return nil
	}

	chatService.Online()

	// The connection will be kept until the client closes it --> call
	// client.CloseHandler.
	<-isAlive
}
