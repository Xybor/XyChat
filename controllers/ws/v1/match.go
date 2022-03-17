package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/helpers/ws"
	"github.com/xybor/xychat/helpers/xybinders"
	resources "github.com/xybor/xychat/resources/v1"
	services "github.com/xybor/xychat/services/v1"
	"github.com/xybor/xychat/xyerrors"
)

// WSMatchHandler handles a connection in websocket protocol and calls match
// service to find a match for client.
func WSMatchHandler(ctx *gin.Context) {
	connection := ctx.MustGet("WebSocket")
	conn := connection.(*websocket.Conn)
	client := ws.CreateWSClient(conn)

	defer client.Close()

	request := resources.WebSocketRequest{}
	xerr := xybinders.Bind(ctx, &request, xybinders.Context)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		client.WriteJSON(response)
		client.Close()
	}

	userService := services.CreateUserService(request.SrcId, true)
	matchService, xerr := services.CreateMatchService(userService)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		client.WriteJSON(response)
		return
	}

	var isAlive = make(chan bool)

	matchService.MatchHandler = func(room resources.RoomResponse) {
		// 0 is an invalid room's identity, therefore it sends failure response to
		// client.
		if room.ID == 0 {
			response := helpers.NewErrorResponse(xyerrors.ErrorUnknown.New("Can't match with anyone"))
			client.WriteJSON(response)
			return
		}

		// It sends room information to the client if success.
		response := helpers.NewResponse(room)
		client.WriteJSON(response)

		isAlive <- false
		close(isAlive)
	}

	// Before closing the wsClient, it needs to unregister the current user
	// from matching queue and close the matchService.
	client.CloseHandler = func() {
		matchService.Unregister()
		matchService.Close()
	}

	// Register the user to matching queue and waiting the result.
	matchService.Register()

	// Keep the connection until the matchQueue returns value, it has two cases:
	// Case 1: The connection closed at client-side --> call client.CloseHandler
	// --> matchService.unregister --> call matchService.MatchHandler -->
	// isAlive signals.
	// Case 2: The matchQueue finds a match --> call matchService.MatchHandler.
	<-isAlive
}
