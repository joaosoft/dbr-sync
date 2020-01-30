package json

import (
	"reflect"
)

func reflectAlloc(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	return reflect.New(typ).Elem()
}
