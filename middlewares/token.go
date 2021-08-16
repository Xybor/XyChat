package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers/context"
	"github.com/xybor/xychat/helpers/tokens"
	"github.com/xybor/xychat/xyerrors"
)

// VerifyUserToken finds the token in incoming request and validates it.  If
// there is a valid token, it will set the token's id as a parameter UID in the
// context;
func VerifyUserToken(cookie bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		var err error
		var xerr xyerrors.XyError

		if !cookie {
			context.SetRetrievingMethod(context.GET)
			token, xerr = context.RetrieveQuery(ctx, "xytok")
			if xerr.Errno() != 0 {
				err = xerr
			} else {
				err = nil
			}
		} else {
			token, err = ctx.Cookie("xytok")
		}

		if err != nil {
			return
		}

		userToken := tokens.CreateEmptyUserToken()

		xerr = userToken.Validate(token)
		if xerr.Errno() != 0 {
			log.Println(xerr)
			return
		}

		id := userToken.GetUID()
		ctx.Set("UID", id)
	}
}
