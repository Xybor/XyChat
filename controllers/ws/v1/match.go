package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ctrl "github.com/xybor/xychat/controllers"
	"github.com/xybor/xychat/helpers/context"
	wshelper "github.com/xybor/xychat/helpers/ws/v1"
	service "github.com/xybor/xychat/services/v1"
)

// MatchHandler handle an incoming request which has already upgraded to
// websocket connection.
//
// It requires the authenticated token to push the user to matching queue.
//
// After it finds a match or occurs an error, it sends a message to the client
// and close the connection.
//
// Note that a user has only a connection to match queue.
func MatchHandler(ctx *gin.Context) {
	connection, _ := ctx.Get("WebSocket")
	conn := connection.(*websocket.Conn)
	client := CreateWSClient(conn)

	defer client.Close()

	id := context.GetUID(ctx)
	userService := service.CreateUserService(id)

	userRepresentation, err := userService.SelfSelect()
	if err != nil {
		response := wshelper.NewWSError(ctrl.ErrorFailedProcess, err.Error())
		client.WriteJSON(response)
		return
	}

	clientService := service.CreateClientService(*userRepresentation)
	// If there was already a connection of this user, the return value is nil.
	if clientService == nil {
		response := wshelper.NewWSError(
			ctrl.ErrorDuplicatedConnection,
			"An account is only allowed to join into the match queue one time",
		)
		client.WriteJSON(response)
		return
	}

	// Before closing the wsClient, it needs to unregister the current user
	// from matching queue and close the service.
	client.CloseHandler = func() {
		defer clientService.Unregister()
		defer clientService.Close()
	}

	// Register the user to matching queue and waiting the result.
	clientService.Register()
	room := clientService.WaitForJoinRoom()

	// 0 is an invalid room's identity, therefore it sends failure response to
	// client.
	if room.ID == 0 {
		response := wshelper.NewWSError(ctrl.ErrorFailedProcess, "Can't match anyone")
		client.WriteJSON(response)
		return
	}

	// It sends room information to the client if success.
	response := wshelper.NewWSResponse(room)
	client.WriteJSON(response)
}
