package dbr

import (
	"fmt"
)

type StmtWhereExists struct {
	stmtSelect *StmtSelect
	isNot      bool

	db *db
}

func newStmtWhereExists(db *db, stmtSelect *StmtSelect, isNot bool) *StmtWhereExists {
	return &StmtWhereExists{
		db:         db,
		stmtSelect: stmtSelect,
		isNot:      isNot,
	}
}

func (stmt *StmtWhereExists) Build() (query string, err error) {
	stmtSelect, err := stmt.stmtSelect.Build()
	if err != nil {
		return "", err
	}

	if stmt.isNot {
		query += fmt.Sprintf("%s ", constFunctionNot)
	}

	return fmt.Sprintf("%s (%s)", constFunctionExists, stmtSelect), nil
}
