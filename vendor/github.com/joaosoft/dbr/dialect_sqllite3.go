package dbr

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type dialectSqlLite3 struct{}

func (d *dialectSqlLite3) Name() string {
	return string(constDialectSqlLite3)
}

func (d *dialectSqlLite3) Encode(value interface{}) string {
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

// https://www.sqlite.org/faq.html
func (d *dialectSqlLite3) EncodeString(value string) string {
	return `'` + strings.Replace(value, `'`, `''`, -1) + `'`
}

// https://www.sqlite.org/lang_expr.html
func (d *dialectSqlLite3) EncodeBool(value bool) string {
	if value {
		return constSqlLite3BoolTrue
	}
	return constSqlLite3BoolFalse
}

// https://www.sqlite.org/lang_datefunc.html
func (d *dialectSqlLite3) EncodeTime(value time.Time) string {
	return `'` + value.UTC().Format(constTimeFormat) + `'`
}

// https://www.sqlite.org/lang_expr.html
func (d *dialectSqlLite3) EncodeBytes(value []byte) string {
	return fmt.Sprintf(`X'%x'`, value)
}

func (d *dialectSqlLite3) EncodeColumn(column interface{}) string {
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

func (d *dialectSqlLite3) Placeholder() string {
	return constSqlLite3PlaceHolder
}
