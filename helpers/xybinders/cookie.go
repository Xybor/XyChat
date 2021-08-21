package xybinders

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/xyerrors"
)

type cookieBinder struct {
}

func (c cookieBinder) Name() string {
	return "cookie"
}

func (c cookieBinder) Bind(ctx *gin.Context, ptr interface{}) xyerrors.XyError {
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get(c.Name())

		if name == "" {
			continue
		}

		value, err := ctx.Cookie(name)
		if err != nil {
			continue
		}

		vField := v.Field(i)
		reflectMap(&vField, value)
	}

	return xyerrors.NoError
}
