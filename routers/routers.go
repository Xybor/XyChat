package routers

import (
	"github.com/gin-gonic/gin"
	apihelpers "github.com/xybor/xychat/helpers/api"
	ctr1 "github.com/xybor/xychat/controllers/api/v1"
	"github.com/xybor/xychat/middlewares"
)

func Route() *gin.Engine {
	router := gin.Default()
	router.Use(apihelpers.ApplyAPIHeader)

	r1 := router.Group("/api/v1")
	r1.Use(middlewares.VerifyUserToken)
	{
		r1.GET("auth", ctr1.AuthenticateUserHandler)
		r1.GET("register", ctr1.RegisterUserHandler)
		r1.GET("profile", ctr1.GetProfileHandler)
	}

	return router
}
