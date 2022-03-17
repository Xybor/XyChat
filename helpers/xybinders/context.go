package xybinders

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/xyerrors"
)

type contextBinder struct {
}

func (c contextBinder) Name() string {
	return "context"
}

func (c contextBinder) Bind(ctx *gin.Context, ptr interface{}) xyerrors.XyError {
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get(c.Name())

		if name == "" {
			continue
		}

		value, ok := ctx.Get(name)
		if !ok {
			continue
		}

		vField := v.Field(i)
		reflectMap(&vField, value)
	}

	return xyerrors.NoError
}
