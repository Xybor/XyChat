package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers/context"
	"github.com/xybor/xychat/helpers/tokens"
)

// VerifyUserToken finds the token in incoming request and validates it.  If
// there is a valid token, it will set the token's id as a parameter UID in the
// context;
func VerifyUserToken(cookie bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		var err error

		if !cookie {
			context.SetRetrievingMethod(context.GET)
			token, err = context.RetrieveQuery(ctx, "xytok")
		} else {
			token, err = ctx.Cookie("xytok")
		}

		if err != nil {
			return
		}

		userToken := tokens.CreateEmptyUserToken()

		err = userToken.Validate(token)
		if err != nil {
			return
		}

		id := userToken.GetUID()
		ctx.Set("UID", id)
	}
}
