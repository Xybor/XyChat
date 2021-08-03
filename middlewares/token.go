package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	ctrl "github.com/xybor/xychat/controllers"
	"github.com/xybor/xychat/helpers/tokens"
)

func VerifyUserToken(c *gin.Context) {
	ctrl.SetReceivingMethod(ctrl.GET)

	// token, err := c.Cookie("auth")
	token, err := ctrl.GetParam(c, "token")

	if err != nil {
		return
	}

	userToken := tokens.CreateUserToken(0, 0)

	err = userToken.Validate(token)
	if err != nil {
		log.Println(err)
		return
	}

	c.Set("UID", userToken.GetUID())
}
