package dbr

import (
	"fmt"
)

type sets struct {
	list []*set
	db   *db
}

func newSets(db *db) *sets {
	return &sets{
		db: db,
		list: make([]*set, 0),
	}
}
func (s sets) Build() (string, error) {

	if len(s.list) == 0 {
		return "", nil
	}

	var query string

	lenS := len(s.list)
	for i, item := range s.list {
		query += fmt.Sprintf("%s = %s", item.column, s.db.Dialect.Encode(item.value))

		if i+1 < lenS {
			query += ", "
		}
	}

	return query, nil
}
