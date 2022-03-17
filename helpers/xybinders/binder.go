package xybinders

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xybor/xychat/xyerrors"
)

type XyBinder interface {
	Bind(*gin.Context, interface{}) xyerrors.XyError
}

var (
	Json    = jsonBinder{}
	Query   = queryBinder{}
	Uri     = uriBinder{}
	Context = contextBinder{}
	Cookie  = cookieBinder{}
)

var validate = validator.New()

// Bind function binds parameters in a gin.Context to the ptr struct.  It
// experiences from binding library and has some more binding types.
//
// The buillt-in binding types are: Json, Query, Uri.
//
// The added binding types are: Context, Cookie.
//
// Please use the binding types in this package instead of original binding
// package if you want to use this function.
//
// Tag: if you want to bind a parameter type, you must declare the
// corresponding tag with the name of that field.  If a field is required, use
// the tag `validate:"required"`.
//
// Example:
//
// type Student struct {
//    ID    uint   `json:"sid" validate:"required"`
//    Name  string `form:"name" context:"name"`
//    Grade int    `form:"g" uri:"grade"`
// }
//
// var student = Student{}
//
// Bind(ctx, &student, xybinders.Json, xybinders.Query)
func Bind(ctx *gin.Context, ptr interface{}, binders ...XyBinder) xyerrors.XyError {
	for _, binder := range binders {
		xerr := binder.Bind(ctx, ptr)
		if xerr.Errno() != 0 {
			return xerr
		}
	}

	err := validate.Struct(ptr)
	if err != nil {
		return xyerrors.ErrorSyntaxInput.New("At least a required field is empty")
	}

	return xyerrors.NoError
}
