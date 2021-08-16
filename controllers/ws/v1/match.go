package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xybor/xychat/helpers/context"
	wshelper "github.com/xybor/xychat/helpers/ws/v1"
	representations "github.com/xybor/xychat/representations/v1"
	service "github.com/xybor/xychat/services/v1"
	xyerrors "github.com/xybor/xychat/xyerrors/v1"
)

// WSMatchHandler handles an incoming request which has already upgraded to
// websocket connection.
//
// It requires the authenticated token to push the user to matchQueue.
//
// After it finds a match or occurs an error, it sends a message to the client
// and closes the connection.
//
// Note that a user has only a connection to the matchQueue.
func WSMatchHandler(ctx *gin.Context) {
	connection := ctx.MustGet("WebSocket")
	conn := connection.(*websocket.Conn)
	client := CreateWSClient(conn)

	defer client.Close()

	id := context.GetUID(ctx)

	userService := service.CreateUserService(id, true)
	matchService, xerr := service.CreateMatchService(userService)
	if xerr.Errno() != 0 {
		response := wshelper.NewWSError(xerr)
		client.WriteJSON(response)
		return
	}

	var isAlive = make(chan bool)

	matchService.MatchHandler = func(room representations.RoomRepresentation) {
		// 0 is an invalid room's identity, therefore it sends failure response to
		// client.
		if room.ID == 0 {
			response := wshelper.NewWSError(xyerrors.ErrorUnknown.New("Can't match with anyone"))
			client.WriteJSON(response)
			return
		}

		// It sends room information to the client if success.
		response := wshelper.NewWSResponse(room)
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
