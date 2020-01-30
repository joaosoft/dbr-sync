package dbr

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type dialectMySql struct{}

func (d *dialectMySql) Name() string {
	return string(constDialectMysql)
}

func (d *dialectMySql) Encode(value interface{}) string {
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

// https://dev.mysql.com/doc/refman/5.7/en/string-literals.html
func (d *dialectMySql) EncodeString(value string) string {
	buf := new(bytes.Buffer)

	buf.WriteRune('\'')
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case 0:
			buf.WriteString(`\0`)
		case '\'':
			buf.WriteString(`\'`)
		case '"':
			buf.WriteString(`\"`)
		case '\b':
			buf.WriteString(`\b`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case 26:
			buf.WriteString(`\Z`)
		case '\\':
			buf.WriteString(`\\`)
		default:
			buf.WriteByte(value[i])
		}
	}

	buf.WriteRune('\'')
	return buf.String()
}

func (d *dialectMySql) EncodeBool(value bool) string {
	if value {
		return constMySqlBoolTrue
	}
	return constMySqlBoolFalse
}

func (d *dialectMySql) EncodeTime(value time.Time) string {
	return `'` + value.UTC().Format(constTimeFormat) + `'`
}

func (d *dialectMySql) EncodeBytes(value []byte) string {
	return fmt.Sprintf(`0x%x`, value)
}

func (d *dialectMySql) EncodeColumn(column interface{}) string {
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

func (d *dialectMySql) Placeholder() string {
	return constMysqlPlaceHolder
}
