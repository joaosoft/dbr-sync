package dbr

import (
	"fmt"
)

type unionType string
type unions []*union

type union struct {
	unionType unionType
	stmt      *StmtSelect
}

func (u unions) Build() (string, error) {

	if len(u) == 0 {
		return "", nil
	}

	var query string

	for _, union := range u {
		stmt, err := union.stmt.Build()
		query += fmt.Sprintf(" %s %s", string(union.unionType), stmt)

		if err != nil {
			return "", err
		}
	}

	return query, nil
}
