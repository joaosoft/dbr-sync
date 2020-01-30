package dbr

import (
	"fmt"
)

type onConflictType string
type onConflictDoType string

const (
	onConflictColumn     onConflictType = "value"
	onConflictConstraint onConflictType = "constraint"

	onConflictDoNothing onConflictDoType = "nothing"
	onConflictDoUpdate  onConflictDoType = "update"
)

type StmtConflict struct {
	onConflictType     onConflictType
	onConflict         *columns
	onConflictDoType   onConflictDoType
	onConflictDoUpdate *sets

	db *db
}

func newStmtConflict(db *db) *StmtConflict {
	return &StmtConflict{
		db:                 db,
		onConflict:         newColumns(db, false),
		onConflictDoUpdate: newSets(db),
	}
}

func (stmt *StmtConflict) Build() (string, error) {

	if stmt.onConflictType == "" {
		return "", nil
	}

	query := fmt.Sprintf(" %s ", constFunctionOnConflict)

	columns, err := stmt.onConflict.Build()
	if err != nil {
		return "", err
	}

	switch stmt.onConflictType {
	case onConflictColumn:
		query += fmt.Sprintf("(%s) ", columns)
	case onConflictConstraint:
		query += fmt.Sprintf("%s (%s) ", constFunctionOnConstraint, columns)
	}

	switch stmt.onConflictDoType {
	case onConflictDoNothing:
		query += constFunctionDoNothing
	case onConflictDoUpdate:
		sets, err := stmt.onConflictDoUpdate.Build()
		if err != nil {
			return "", err
		}
		query += fmt.Sprintf("%s %s", constFunctionDoUpdateSet, sets)
	}

	return query, nil
}
