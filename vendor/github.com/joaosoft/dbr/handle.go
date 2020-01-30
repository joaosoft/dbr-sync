package dbr

import (
	"fmt"
	"reflect"
)

func handleExpression(base *functionBase, expression interface{}) (string, error) {
	var value string
	var err error

	if expression == nil || (reflect.ValueOf(expression).Kind() == reflect.Ptr && reflect.ValueOf(expression).IsNil()) {
		value = fmt.Sprintf(constFunctionNull)
		return value, nil
	}

	switch stmt := expression.(type) {
	case *StmtSelect:
		value, err = stmt.Build()
		if err != nil {
			return "", err
		}
		value = fmt.Sprintf("(%s)", value)
	case Builder:
		value, err = stmt.Build()
		if err != nil {
			return "", nil
		}
	case functionBuilder:
		var err error
		value, err = stmt.Build(base.db)
		if err != nil {
			return "", nil
		}
	case iFunction:
		var err error
		value, err = stmt.Expression(base.db)
		if err != nil {
			return "", nil
		}
	default:
		value = fmt.Sprintf("%+v", expression)
	}

	return value, nil
}

func handleBuild(base *functionBase, expression interface{}, encode ...bool) (string, error) {
	var value string
	var err error

	theEncode := base.encode

	if len(encode) > 0 {
		theEncode = encode[0]
	}

	if expression == nil || (reflect.ValueOf(expression).Kind() == reflect.Ptr && reflect.ValueOf(expression).IsNil()) {
		value = fmt.Sprintf(constFunctionNull)
		return value, nil
	}

	switch stmt := expression.(type) {
	case *StmtSelect:
		value, err = stmt.Build()
		if err != nil {
			return "", err
		}
		value = fmt.Sprintf("(%s)", value)
	case Builder:
		value, err = stmt.Build()
		if err != nil {
			return "", nil
		}
	case functionBuilder:
		var err error
		value, err = stmt.Build(base.db)
		if err != nil {
			return "", nil
		}
	case iFunction:
		var err error
		value, err = stmt.Build(base.db)
		if err != nil {
			return "", nil
		}
	default:
		if theEncode {
			if base.isColumn {
				value = base.db.Dialect.EncodeColumn(expression)
			} else {
				value = base.db.Dialect.Encode(expression)
			}
		} else {
			value = fmt.Sprintf("%+v", expression)
		}
	}

	return value, nil
}
