package dbr

import (
	"database/sql/driver"
	"fmt"
)

type values struct {
	list []interface{}
	*functionBase
}

func newValues(db *db) *values {
	return &values{
		functionBase: newFunctionBase(true, false, db),
	}
}

func (v values) Build() (query string, err error) {
	lenV := len(v.list)
	var withoutParentheses bool

	for i, item := range v.list {
		var value string

		switch stmt := item.(type) {
		case *values:
			withoutParentheses = true
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
		case driver.Valuer:
			valuer, err := stmt.Value()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("%+v", valuer)
		default:
			value, err = handleBuild(v.functionBase, item)
			if err != nil {
				return "", err
			}
		}

		query += value

		if i+1 < lenV {
			query += ", "
		}
	}

	if withoutParentheses {
		return query, nil
	}

	return fmt.Sprintf("(%s)", query), nil
}
