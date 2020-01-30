package dbr

import (
	"database/sql"
	"strings"
	"time"
)

type StmtExecute struct {
	query  string
	values *values

	Dbr          *Dbr
	Db           *db
	Duration     time.Duration
	sqlOperation SqlOperation
}

func newStmtExecute(dbr *Dbr, db *db, query string) *StmtExecute {
	return &StmtExecute{
		sqlOperation: ExecuteOperation,
		Dbr:          dbr,
		Db:           db,
		query:        query,
		values:       newValues(dbr.Connections.Write),
	}
}

func (stmt *StmtExecute) Values(valuesList ...interface{}) *StmtExecute {
	stmt.values.list = append(stmt.values.list, valuesList...)
	return stmt
}

func (stmt *StmtExecute) Build() (query string, _ error) {
	query = stmt.query

	if strings.Count(query, stmt.Db.Dialect.Placeholder()) != len(stmt.values.list) {
		return "", ErrorNumberOfConditionValues
	}

	for _, value := range stmt.values.list {
		query = strings.Replace(query, stmt.Db.Dialect.Placeholder(), stmt.Db.Dialect.Encode(value), 1)
	}

	return query, nil
}

func (stmt *StmtExecute) Exec() (sql.Result, error) {

	startTime := time.Now()
	defer func() {
		stmt.Duration = time.Since(startTime)
	}()

	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	result, err := stmt.Db.Exec(query)
	if err != nil {
		return nil, err
	}

	if err := stmt.Dbr.eventHandler(stmt.sqlOperation, []string{}, query, err, nil, result); err != nil {
		return nil, err
	}

	return result, err
}
