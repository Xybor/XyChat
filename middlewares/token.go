package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers/context"
	"github.com/xybor/xychat/helpers/tokens"
)

// VerifyUserToken finds the token in incoming request and validates it.  If
// there is a valid token, it will set the token's id as a parameter UID in the
// context;
func VerifyUserToken(ctx *gin.Context) {
	context.SetRetrievingMethod(context.GET)

	// token, err := c.Cookie("auth")
	token, err := context.RetrieveQuery(ctx, "token")

	if err != nil {
		return
	}

	userToken := tokens.CreateEmptyUserToken()

	err = userToken.Validate(token)
	if err != nil {
		log.Println(err)
		return
	}

	id := userToken.GetUID()
	ctx.Set("UID", id)
}
