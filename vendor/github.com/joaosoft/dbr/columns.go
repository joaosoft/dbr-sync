package dbr

import (
	"fmt"
)

type columns struct {
	list   []interface{}
	encode bool

	*functionBase
}

func newColumns(db *db, encode bool) *columns {
	return &columns{
		functionBase: newFunctionBase(encode, true, db),
		list:         make([]interface{}, 0),
		encode:       encode,
	}
}

func (c columns) Build() (query string, err error) {

	lenC := len(c.list)
	for i, item := range c.list {
		var value string

		switch stmt := item.(type) {
		case *StmtSelect:
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("(%s)", value)

		default:
			value, err = handleBuild(c.functionBase, item)
			if err != nil {
				return "", err
			}
		}

		query += value

		if i+1 < lenC {
			query += ", "
		}
	}

	return query, nil
}
