package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api1 "github.com/xybor/xychat/controllers/api/v1"
	ws1 "github.com/xybor/xychat/controllers/ws/v1"
	apihelpers "github.com/xybor/xychat/helpers/api"
	"github.com/xybor/xychat/middlewares"
)

func Route() *gin.Engine {
	router := gin.Default()

	router.StaticFS("/ui", http.Dir("vue/dist"))

	rapi := router.Group("/api")
	rapi.Use(
		apihelpers.ApplyAPIHeader,
		middlewares.VerifyUserToken,
	)
	{
		rapi1 := rapi.Group("/v1")
		{
			rapi1.GET("auth", api1.AuthenticateUserHandler)
			rapi1.GET("register", api1.RegisterUserHandler)
			rapi1.GET("profile", api1.GetProfileHandler)
		}
	}
	rws := router.Group("/ws")
	rws.Use(
		middlewares.VerifyUserToken,
		middlewares.UpgradeToWebSocket,
	)
	{
		rws1 := rws.Group("/v1")
		{
			rws1.Any("match", ws1.MatchHandler)
		}
	}
	return router
}
