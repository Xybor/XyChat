package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api1 "github.com/xybor/xychat/controllers/api/v1"
	ws1 "github.com/xybor/xychat/controllers/ws/v1"
	"github.com/xybor/xychat/helpers/context"
	"github.com/xybor/xychat/middlewares"
	mdwv1 "github.com/xybor/xychat/middlewares/v1"
)

// Route combines middlewares and controllers to handle given url paths in the
// application.
func Route() *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.ApplyCORSHeader)

	router.StaticFS("/ui", http.Dir("vue/dist"))

	rapi := router.Group("api")
	rapi.Use(
		middlewares.VerifyUserToken(true),
		middlewares.ApplyAPIHeader,
	)
	{
		rapi1 := rapi.Group("v1")
		{
			rapi1.POST("auth",
				mdwv1.MustHaveQueryParam(context.POST, "username", "password"),
				api1.UserAuthenticateHandler,
			)
			rapi1.POST("register",
				mdwv1.MustHaveQueryParam(context.POST, "username", "password"),
				api1.UserRegisterHandler,
			)
			rapi1.GET("profile", api1.UserProfileHandler)

			rapi1.GET("users/:id", api1.UserGETHandler)
			rapi1.PUT("users/:id", api1.UserPUTHandler)
			rapi1.PUT("users/:id/role",
				mdwv1.MustHaveQueryParam(context.POST, "role"),
				api1.UserChangeRoleHandler,
			)
			rapi1.PUT("users/:id/password",
				mdwv1.MustHaveQueryParam(context.POST, "newpassword"),
				api1.UserChangePasswordHandler,
			)
		}
	}

	rws := router.Group("ws")
	rws.Use(
		middlewares.VerifyUserToken(false),
		middlewares.UpgradeToWebSocket,
	)
	{
		rws1 := rws.Group("v1")
		{
			rws1.GET("match", ws1.MatchHandler)
		}
	}

	return router
}
