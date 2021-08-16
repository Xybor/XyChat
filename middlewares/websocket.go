package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	wshelpers "github.com/xybor/xychat/helpers/ws"
	"github.com/xybor/xychat/xyerrors"
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
		response := wshelpers.NewWSError(
			xyerrors.ErrorCannotUpgradeToWebsocket.New("Websocket is not allowed"))
		ctx.JSON(xyerrors.ErrorCannotUpgradeToWebsocket.StatusCode(), response)
		ctx.Abort()
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		response := wshelpers.NewWSError(
			xyerrors.ErrorUnknown.New("Can't upgrade to websocket because unknown reason"))
		ctx.JSON(xyerrors.ErrorUnknown.StatusCode(), response)
		ctx.Abort()
		return
	}

	ctx.Set("WebSocket", conn)
}
