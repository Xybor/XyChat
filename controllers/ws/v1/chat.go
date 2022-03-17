package v1

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/helpers/ws"
	"github.com/xybor/xychat/helpers/xybinders"
	resources "github.com/xybor/xychat/resources/v1"
	services "github.com/xybor/xychat/services/v1"
	"github.com/xybor/xychat/xyerrors"
)

// WSChatHandler handles a connection in a websocket protocol and communicates
// chat messages with client.
func WSChatHandler(ctx *gin.Context) {
	connection := ctx.MustGet("WebSocket")
	conn := connection.(*websocket.Conn)
	client := ws.CreateWSClient(conn)

	request := resources.WebSocketRequest{}
	xerr := xybinders.Bind(ctx, &request, xybinders.Context)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		client.WriteJSON(response)
		client.Close()
	}

	userService := services.CreateUserService(request.SrcId, true)
	chatService, xerr := services.CreateChatService(userService)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
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
		msg := resources.ChatMessageRequest{}
		err := json.Unmarshal(message, &msg)
		if err != nil {
			response := helpers.NewErrorResponse(xyerrors.ErrorSyntaxInput.New("invalid json format"))
			client.WriteJSON(response)
			return nil
		}

		xerr = chatService.SendTo(msg.RoomId, msg.Message)
		if xerr.Errno() != 0 {
			response := helpers.NewErrorResponse(xerr)
			client.WriteJSON(response)
			return nil
		}

		return nil
	}

	// Forward the message from broadcast service to client
	chatService.ChatHandler = func(cmr resources.ChatMessageResponse) error {
		response := helpers.NewResponse(cmr)
		client.WriteJSON(response)
		return nil
	}

	chatService.Online()

	// The connection will be kept until the client closes it --> call
	// client.CloseHandler.
	<-isAlive
}
