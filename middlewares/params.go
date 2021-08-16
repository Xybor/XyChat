package middlewares

import (
	"github.com/gin-gonic/gin"
	apihelper "github.com/xybor/xychat/helpers/api"
	"github.com/xybor/xychat/helpers/context"
)

// MustHaveQueryParam is a middleware checking if the provided parameters exist
// in the request or not.  The method indicates the position to find
// parameters.
//
// If any parameter doesn't exist, it will respond a http BadRequest
// immediately.
func MustHaveQueryParam(method string, params ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context.SetRetrievingMethod(method)

		for _, param := range params {
			_, xerr := context.RetrieveQuery(ctx, param)

			if xerr.Errno() != 0 {
				response := apihelper.NewAPIError(xerr)
				ctx.JSON(xerr.StatusCode(), response)
				ctx.Abort()
				return
			}
		}
	}
}
