package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

type StmtInsert struct {
	withStmt     *StmtWith
	table        *table
	columns      *columns
	values       *values
	returning    *columns
	stmtConflict *StmtConflict
	fromSelect   *StmtSelect

	Dbr          *Dbr
	Db           *db
	Duration     time.Duration
	sqlOperation SqlOperation
}

func newStmtInsert(dbr *Dbr, db *db, withStmt *StmtWith) *StmtInsert {
	return &StmtInsert{
		sqlOperation: InsertOperation,
		Dbr:          dbr,
		Db:           db,
		withStmt:     withStmt,
		values:       newValues(dbr.Connections.Write),
		stmtConflict: newStmtConflict(dbr.Connections.Write),
		columns:      newColumns(dbr.Connections.Write, false),
		returning:    newColumns(dbr.Connections.Write, false),
	}
}

func (stmt *StmtInsert) Into(table interface{}) *StmtInsert {
	stmt.table = newTable(stmt.Db, table)
	return stmt
}

func (stmt *StmtInsert) Columns(columns ...interface{}) *StmtInsert {
	stmt.columns.list = append(stmt.columns.list, columns...)
	return stmt
}

func (stmt *StmtInsert) Values(valuesList ...interface{}) *StmtInsert {
	stmt.values.list = append(stmt.values.list, &values{functionBase: newFunctionBase(true, false, stmt.Db), list: valuesList})
	return stmt
}

func (stmt *StmtInsert) FromSelect(selectStmt *StmtSelect) *StmtInsert {
	stmt.fromSelect = selectStmt
	return stmt
}

func (stmt *StmtInsert) Build() (string, error) {
	var query string

	// withStmt
	if len(stmt.withStmt.withs) > 0 {
		withStmt, err := stmt.withStmt.Build()
		if err != nil {
			return "", err
		}
		query += withStmt
	}

	// columns
	columns, err := stmt.columns.Build()
	if err != nil {
		return "", err
	}

	// values
	var values string
	if len(stmt.values.list) > 0 {
		values = fmt.Sprintf("%s ", constFunctionValues)
		val, err := stmt.values.Build()
		if err != nil {
			return "", err
		}
		values += val
	}

	// from select statement
	var selectStmt string
	if stmt.fromSelect != nil {
		selectStmt, err = stmt.fromSelect.Build()
		if err != nil {
			return "", err
		}
	}

	table, err := stmt.table.Build()
	if err != nil {
		return "", err
	}

	query += fmt.Sprintf("%s %s %s (%s) %s%s", constFunctionInsert, constFunctionInto, table, columns, values, selectStmt)

	// on conflict
	if stmt.stmtConflict.onConflictType != "" {
		onConflictStmt, err := stmt.stmtConflict.Build()
		if err != nil {
			return "", err
		}

		query += onConflictStmt
	}

	// returning
	if len(stmt.returning.list) > 0 {
		returning, err := stmt.returning.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" %s %s", constFunctionReturning, returning)
	}

	return query, nil
}

func (stmt *StmtInsert) Exec() (sql.Result, error) {

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

func (stmt *StmtInsert) Record(record interface{}) *StmtInsert {
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

	var valueList []interface{}
	for _, column := range stmt.columns.list {
		valueList = append(valueList, mappedValues[column].Interface())
	}

	stmt.values.list = append(stmt.values.list, &values{functionBase: newFunctionBase(true, false, stmt.Db), list: valueList})

	return stmt
}

func (stmt *StmtInsert) Records(records []interface{}) *StmtInsert {
	for _, record := range records {
		stmt.Record(record)
	}

	return stmt
}

func (stmt *StmtInsert) OnConflict(column ...interface{}) *StmtInsert {
	stmt.stmtConflict.onConflictType = onConflictColumn
	stmt.stmtConflict.onConflict.list = append(stmt.stmtConflict.onConflict.list, column...)
	return stmt
}

func (stmt *StmtInsert) OnConflictConstraint(constraint interface{}) *StmtInsert {
	stmt.stmtConflict.onConflictType = onConflictConstraint
	stmt.stmtConflict.onConflict.list = []interface{}{constraint}
	return stmt
}

func (stmt *StmtInsert) DoNothing() *StmtInsert {
	stmt.stmtConflict.onConflictDoType = onConflictDoNothing
	return stmt
}

func (stmt *StmtInsert) DoUpdate(fieldValue ...interface{}) *StmtInsert {
	stmt.stmtConflict.onConflictDoType = onConflictDoUpdate

	if len(fieldValue)%2 != 0 {
		return stmt
	}

	lenC := len(fieldValue)
	for i := 0; i < lenC; i += 2 {
		stmt.stmtConflict.onConflictDoUpdate.list = append(stmt.stmtConflict.onConflictDoUpdate.list, &set{column: fieldValue[i].(string), value: fieldValue[i+1]})
	}

	return stmt
}

func (stmt *StmtInsert) Return(column ...interface{}) *StmtInsert {
	stmt.returning.list = append(stmt.returning.list, column...)
	return stmt
}

func (stmt *StmtInsert) Load(object interface{}) (count int, err error) {

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
