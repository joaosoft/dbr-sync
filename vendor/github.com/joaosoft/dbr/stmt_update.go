package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

type StmtUpdate struct {
	withStmt   *StmtWith
	table      *table
	sets       *sets
	columns    *columns
	conditions *conditions
	returning  *columns

	Dbr          *Dbr
	Db           *db
	Duration     time.Duration
	sqlOperation SqlOperation
}

func newStmtUpdate(dbr *Dbr, db *db, withStmt *StmtWith, table interface{}) *StmtUpdate {
	return &StmtUpdate{
		sqlOperation: UpdateOperation,
		Dbr:          dbr,
		Db:           db,
		withStmt:     withStmt,
		table:        newTable(db, table),
		sets:         newSets(dbr.Connections.Write),
		conditions:   newConditions(dbr.Connections.Write),
		columns:      newColumns(dbr.Connections.Write, false),
		returning:    newColumns(dbr.Connections.Write, false),
	}
}

func (stmt *StmtUpdate) Set(column string, value interface{}) *StmtUpdate {
	stmt.sets.list = append(stmt.sets.list, &set{column: column, value: value})
	return stmt
}

func (stmt *StmtUpdate) From(table interface{}) *StmtUpdate {
	stmt.table = newTable(stmt.Db, table)
	return stmt
}

func (stmt *StmtUpdate) Where(query interface{}, values ...interface{}) *StmtUpdate {
	stmt.conditions.list = append(stmt.conditions.list, &condition{operator: OperatorAnd, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtUpdate) WhereOr(query string, values ...interface{}) *StmtUpdate {
	stmt.conditions.list = append(stmt.conditions.list, &condition{operator: OperatorOr, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtUpdate) Build() (query string, err error) {
	// withStmt
	withStmt, err := stmt.withStmt.Build()
	if err != nil {
		return "", err
	}
	query += withStmt

	sets, err := stmt.sets.Build()
	if err != nil {
		return "", err
	}

	table, err := stmt.table.Build()
	if err != nil {
		return "", err
	}

	query += fmt.Sprintf("%s %s %s %s",  constFunctionUpdate, table, constFunctionSet, sets)

	if len(stmt.conditions.list) > 0 {
		conds, err := stmt.conditions.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" %s %s", constFunctionWhere, conds)
	}

	if len(stmt.returning.list) > 0 {
		returning, err := stmt.returning.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" %s %s", constFunctionReturning, returning)
	}

	return query, nil
}

func (stmt *StmtUpdate) Exec() (sql.Result, error) {

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

	if err := stmt.Dbr.eventHandler(stmt.sqlOperation, []string{fmt.Sprint(stmt.table)}, query, err, nil, result); err != nil {
		return nil, err
	}

	return result, err
}

func (stmt *StmtUpdate) Record(record interface{}) *StmtUpdate {
	value := reflect.ValueOf(record)

	mappedValues := make(map[interface{}]reflect.Value)

	if len(stmt.columns.list) == 0 {
		var columns []interface{}
		loadStructValues(constLoadOptionWrite, value, &columns, mappedValues)
		stmt.columns.list = columns
		stmt.columns.encode = true
	} else {
		loadStructValues(constLoadOptionWrite, value, nil, mappedValues)
	}

	for _, column := range stmt.columns.list {
		stmt.sets.list = append(stmt.sets.list, &set{column: column, value: mappedValues[column].Interface()})
	}

	return stmt
}

func (stmt *StmtUpdate) Return(column ...interface{}) *StmtUpdate {
	stmt.returning.list = append(stmt.returning.list, column...)
	return stmt
}

func (stmt *StmtUpdate) Load(object interface{}) (count int, err error) {

	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return 0, ErrorInvalidPointer
	}

	startTime := time.Now()
	defer func() {
		stmt.Duration = time.Since(startTime)
	}()

	query, err := stmt.Build()
	if err != nil {
		return 0, err
	}

	rows, err := stmt.Db.Query(query)
	if err != nil {
		return 0, err
	}

	if err := stmt.Dbr.eventHandler(stmt.sqlOperation, []string{fmt.Sprint(stmt.table)}, query, err, rows, nil); err != nil {
		return 0, err
	}

	defer rows.Close()

	return read(stmt.returning.list, rows, value)
}
