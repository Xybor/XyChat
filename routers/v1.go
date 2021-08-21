package routers

import (
	"github.com/gin-gonic/gin"
	api1 "github.com/xybor/xychat/controllers/api/v1"
	ws1 "github.com/xybor/xychat/controllers/ws/v1"
)

func routeAPIv1(rapi1 *gin.RouterGroup) {
	rapi1.POST("auth",
		api1.UserAuthenticateHandler,
	)
	rapi1.POST("register",
		api1.UserRegisterHandler,
	)
	rapi1.GET("profile", api1.UserProfileHandler)

	rapi1.GET("users/:id", api1.UserSelectHandler)
	rapi1.PUT("users/:id", api1.UserUpdateHandler)
	rapi1.PUT("users/:id/role", api1.UserUpdateRoleHandler)
	rapi1.PUT("users/:id/password", api1.UserChangePasswordHandler)
}

func routeWSv1(rws1 *gin.RouterGroup) {
	rws1.GET("match", ws1.WSMatchHandler)
	rws1.GET("chat", ws1.WSChatHandler)
}
