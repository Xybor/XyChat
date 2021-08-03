package v1

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ctrl "github.com/xybor/xychat/controllers"
	wshelper "github.com/xybor/xychat/helpers/ws/v1"
	service "github.com/xybor/xychat/services/v1"
)

func MatchHandler(ctx *gin.Context) {
	connection, _ := ctx.Get("WebSocket")
	conn := connection.(*websocket.Conn)
	client := CreateWSClient(conn)

	defer client.Close()

	suid, ok := ctx.Get("UID")
	if !ok {
		response := wshelper.NewWSError(ctrl.Unauthenticated, "Unauthenticated")
		client.WriteJSON(response)
		return
	}

	uid := suid.(uint)
	userService := service.CreateUserService(&uid, nil, nil)
	userRepresentation, err := userService.GetProfile()
	if err != nil {
		log.Println(err)
		response := wshelper.NewWSError(ctrl.FailedProcess, "Invalid uid")
		client.WriteJSON(response)
		return
	}

	clientService := service.CreateClientService(userRepresentation)
	if clientService == nil {
		response := wshelper.NewWSError(
			ctrl.DuplicatedConnection,
			"An account is only allowed to join into the match queue one time",
		)
		client.WriteJSON(response)
		return
	}

	client.CloseHandler = func() {
		defer clientService.Unregister()
		defer clientService.Close()
	}

	clientService.Register()
	room := clientService.WaitForJoinRoom()

	if room.ID == 0 {
		response := wshelper.NewWSError(ctrl.FailedProcess, "Can't match anyone")
		client.WriteJSON(response)
		return
	}

	response := wshelper.NewWSResponse(room)
	client.WriteJSON(response)
}
