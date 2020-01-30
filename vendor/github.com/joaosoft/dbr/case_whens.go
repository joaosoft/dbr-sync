package dbr

import (
	"fmt"
)

type onCaseWhens []*caseWhen

func newCaseWhens() onCaseWhens {
	return make(onCaseWhens, 0)
}

func (c onCaseWhens) Build(db *db) (query string, err error) {
	if len(c) == 0 {
		return "", nil
	}

	for i, cond := range c {
		onWhen, err := cond.Build(db)
		if err != nil {
			return "", err
		}

		if i > 0 {
			query += " "
		}

		query += fmt.Sprintf("%s", onWhen)
	}

	return query, nil
}
