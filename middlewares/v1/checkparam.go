package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ctrl "github.com/xybor/xychat/controllers"
	apihelper "github.com/xybor/xychat/helpers/api/v1"
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
			_, err := context.RetrieveQuery(ctx, param)

			if err != nil {
				response := apihelper.NewAPIError(ctrl.ErrorLackOfInput, "empty "+param)
				ctx.JSON(http.StatusBadRequest, response)
				ctx.Abort()
				return
			}
		}
	}
}
