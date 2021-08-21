package xybinders

import (
	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/xyerrors"
)

type uriBinder struct {
}

func (u uriBinder) Bind(ctx *gin.Context, ptr interface{}) xyerrors.XyError {
	if err := ctx.ShouldBindUri(ptr); err != nil {
		return xyerrors.ErrorSyntaxInput.New("Invalid syntax")
	}

	return xyerrors.NoError
}
