package xybinders

import (
	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/xyerrors"
)

type queryBinder struct {
}

func (q queryBinder) Bind(ctx *gin.Context, ptr interface{}) xyerrors.XyError {
	if err := ctx.ShouldBindQuery(ptr); err != nil {
		return xyerrors.ErrorSyntaxInput.New("Invalid syntax")
	}

	return xyerrors.NoError
}
