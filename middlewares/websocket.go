package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xybor/xychat/controllers"
	wshelper "github.com/xybor/xychat/helpers/ws/v1"
)

var upgrader = &websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
}

// UpgradeToWebSocket upgrades the current HTTP connection to WebSocket
// connection and sets it as the parameter WebSocket in the context.  If it
// cannot upgrade to WebSocket connection, a failure response will be sent to
// client.
func UpgradeToWebSocket(ctx *gin.Context) {
	if !ctx.IsWebsocket() {
		response := wshelper.NewWSError(
			controllers.ErrorFailedProcess,
			"websocket connection is required",
		)
		ctx.JSON(http.StatusMethodNotAllowed, response)
		ctx.Abort()
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		response := wshelper.NewWSError(
			controllers.ErrorFailedProcess,
			"cannot upgrade to websocket",
		)
		ctx.JSON(http.StatusInternalServerError, response)
		ctx.Abort()
		return
	}

	ctx.Set("WebSocket", conn)
}
