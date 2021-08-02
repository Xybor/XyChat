package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	ctr "github.com/xybor/xychat/controllers"
	"github.com/xybor/xychat/helpers/tokens"
)

func VerifyUserToken(c *gin.Context) {
	ctr.SetReceivingMethod(ctr.GET)

	// token, err := c.Cookie("auth")
	token, err := ctr.GetParam(c, "token")

	if err != nil {
		return
	}

	ut := tokens.UserToken{}

	err = ut.Validate(token)
	if err != nil {
		log.Println(err)
		return
	}

	c.Set("UID", ut.ID)
}
