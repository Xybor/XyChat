package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ctr1 "github.com/xybor/xychat/controllers/api/v1"
	apihelpers "github.com/xybor/xychat/helpers/api"
	"github.com/xybor/xychat/middlewares"
)

func Route() *gin.Engine {
	router := gin.Default()

	router.StaticFS("/ui", http.Dir("vue/dist"))

	r1 := router.Group("/api/v1")
	r1.Use(apihelpers.ApplyAPIHeader)
	r1.Use(middlewares.VerifyUserToken)
	{
		r1.GET("auth", ctr1.AuthenticateUserHandler)
		r1.GET("register", ctr1.RegisterUserHandler)
		r1.GET("profile", ctr1.GetProfileHandler)
	}

	return router
}
