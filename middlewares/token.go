package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers/tokens"
	"github.com/xybor/xychat/helpers/xybinders"
	"github.com/xybor/xychat/xyerrors"
)

type XyTokRequest struct {
	Xytok string `form:"xytok" cookie:"xytok" validate:"required"`
}

// VerifyUserToken finds the token in incoming request and validates it.  If
// there is a valid token, it will set the token's id as a parameter UID in the
// context;
func VerifyUserToken(cookie bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var xerr xyerrors.XyError
		request := XyTokRequest{}

		if !cookie {
			xerr = xybinders.Bind(ctx, &request, xybinders.Query)
		} else {
			xerr = xybinders.Bind(ctx, &request, xybinders.Cookie)
		}

		if xerr.Errno() != 0 {
			return
		}

		userToken := tokens.CreateEmptyUserToken()

		xerr = userToken.Validate(request.Xytok)
		if xerr.Errno() != 0 {
			log.Println(xerr)
			return
		}

		id := userToken.GetUID()
		ctx.Set("id", id)
	}
}
