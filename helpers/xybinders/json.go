package xybinders

import (
	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/xyerrors"
)

type jsonBinder struct {
}

func (j jsonBinder) Bind(ctx *gin.Context, ptr interface{}) xyerrors.XyError {
	if err := ctx.ShouldBindJSON(ptr); err != nil {
		return xyerrors.ErrorSyntaxInput.New("Invalid syntax")
	}

	return xyerrors.NoError
}
