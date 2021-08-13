package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xybor/xychat/controllers"
	"github.com/xybor/xychat/helpers/context"
	wshelpers "github.com/xybor/xychat/helpers/ws/v1"
	r "github.com/xybor/xychat/representations/v1"
	services "github.com/xybor/xychat/services/v1"
)

type ClientMessage struct {
	RoomId  uint   `json:"roomid"`
	Message string `json:"message"`
}

func WSChatHandler(ctx *gin.Context) {
	connection, _ := ctx.Get("WebSocket")
	conn := connection.(*websocket.Conn)
	client := CreateWSClient(conn)

	uid := context.GetUID(ctx)
	userService := services.CreateUserService(uid)

	chatService, err := services.CreateChatService(userService)

	if err != nil {
		response := wshelpers.NewWSError(controllers.ErrorUnauthenticated, err.Error())
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	var isAlive = make(chan bool)

	client.CloseHandler = func() {
		chatService.Offline()

		client.Close()

		isAlive <- false
		close(isAlive)
	}

	client.ReadHandler = func(message []byte) error {
		msg := ClientMessage{}
		err := json.Unmarshal(message, &msg)
		if err != nil {
			response := wshelpers.NewWSError(controllers.ErrorInvalidInput, "invalid json format")
			client.WriteJSON(response)
			return nil
		}

		err = chatService.SendTo(msg.RoomId, msg.Message)
		if err != nil {
			response := wshelpers.NewWSError(controllers.ErrorFailedProcess, err.Error())
			client.WriteJSON(response)
			return nil
		}

		return nil
	}

	chatService.ChatHandler = func(cmr r.ChatMessageRepresentation) error {
		response := wshelpers.NewWSResponse(cmr)
		client.WriteJSON(response)
		return nil
	}

	chatService.Online()

	<-isAlive
}
