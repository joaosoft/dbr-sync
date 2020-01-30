package dbr

import (
	"fmt"
)

type joins []*StmtJoin

func (j joins) Build() (string, error) {

	if len(j) == 0 {
		return "", nil
	}

	var query string

	lenJ := len(j)
	for i, item := range j {
		join, err := item.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf("%s", join)

		if i+1 < lenJ {
			query += " "
		}
	}

	return query, nil
}
