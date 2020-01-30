package dbr

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type dialectPostgres struct{}

func (d *dialectPostgres) Name() string {
	return string(constDialectPostgres)
}

func (d *dialectPostgres) Encode(value interface{}) string {
	encoded := reflect.ValueOf(value)

	if encoded.Kind() == reflect.Ptr {
		if encoded.IsNil() {
			return constFunctionNull
		}
		encoded = encoded.Elem()
	}

	switch encoded.Kind() {
	case reflect.String:
		return d.EncodeString(encoded.String())
	case reflect.Bool:
		return d.EncodeBool(encoded.Bool())
	default:
		switch encoded.Type() {
		case reflect.TypeOf(time.Time{}):
			return d.EncodeTime(value.(time.Time))
		case reflect.TypeOf([]byte{}):
			return d.EncodeBytes(value.([]byte))
		}
	}

	return fmt.Sprintf("%+v", encoded.Interface())
}

func (d *dialectPostgres) EncodeString(value string) string {
	return `'` + strings.Replace(value, `'`, `''`, -1) + `'`
}

func (d *dialectPostgres) EncodeBool(value bool) string {
	if value {
		return constPostgresBoolTrue
	}
	return constPostgresBoolFalse
}

func (d *dialectPostgres) EncodeTime(value time.Time) string {
	return `'` + value.UTC().Format(constTimeFormat) + `'`
}

func (d *dialectPostgres) EncodeBytes(value []byte) string {
	return fmt.Sprintf(`E'\\x%x'`, value)
}

func (d *dialectPostgres) EncodeColumn(column interface{}) string {
	value := fmt.Sprintf("%+v", column)

	switch column.(type) {
	case string:
		if !strings.ContainsAny(value, `*`) {
			value = fmt.Sprintf(`"%s"`, value)
			value = strings.Replace(value, `.`, `"."`, 1)
		}
	}

	return value
}

func (d *dialectPostgres) Placeholder() string {
	return constPostgresPlaceHolder
}
