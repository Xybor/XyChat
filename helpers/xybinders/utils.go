package xybinders

import (
	"log"
	"reflect"
)

func reflectMap(vField *reflect.Value, value interface{}) {
	vFieldType := vField.Kind()
	vFieldActualType := vFieldType

	if vFieldType == reflect.Ptr {
		vFieldActualType = vField.Type().Elem().Kind()
	}

	switch vFieldActualType {
	case reflect.Uint:
		reflectUintMap(vField, value)
	case reflect.String:
		reflectStringMap(vField, value)
	default:
		log.Panicf("Unknown type %s", vField.Kind())
	}
}

func reflectUintMap(vField *reflect.Value, value interface{}) {
	valuecast := value.(uint)
	if vField.Kind() == reflect.Ptr {
		vField.Set(reflect.ValueOf(&valuecast))
	} else {
		vField.SetUint(uint64(valuecast))
	}
}

func reflectStringMap(vField *reflect.Value, value interface{}) {
	valuecast := value.(string)
	if vField.Kind() == reflect.Ptr {
		vField.Set(reflect.ValueOf(&valuecast))
	} else {
		vField.SetString(valuecast)
	}
}
