package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/middlewares"
)

// Route combines middlewares and controllers to handle given urls in the
// application.
func Route() *gin.Engine {
	router := gin.Default()
	rapi := router.Group("api")
	rapi.Use(
		middlewares.VerifyUserToken(true),
		middlewares.ApplyCORSHeader(),
	)
	{
		api1 := rapi.Group("v1")
		routeAPIv1(api1)
	}

	rws := router.Group("ws")
	rws.Use(
		middlewares.VerifyUserToken(false),
		middlewares.UpgradeToWebSocket(),
	)
	{
		rws1 := rws.Group("v1")
		routeWSv1(rws1)
	}
	return router
}
