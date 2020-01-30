package dbr

import (
	"fmt"
)

type direction string

type order struct {
	column    string
	direction direction
}

type orders []*order

func (o orders) Build() (query string, _ error) {
	if len(o) == 0 {
		return "", nil
	}

	query = fmt.Sprintf(" %s ", constFunctionOrderBy)

	lenO := len(o)
	for i, item := range o {
		query += fmt.Sprintf("%s %s", item.column, item.direction)

		if i+1 < lenO {
			query += ", "
		}
	}

	return query, nil
}
