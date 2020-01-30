package dbr

import (
	"fmt"
	"reflect"
	"time"
)

type StmtSelect struct {
	withStmt          *StmtWith
	columns           *columns
	tables            tables
	joins             joins
	existsStmt        *StmtWhereExists
	conditions        *conditions
	isDistinct        bool
	distinctColumns   *columns
	distinctOnColumns *columns
	unions            unions
	groupBy           groupBy
	having            *conditions
	orders            orders
	returning         *columns
	limit             int
	offset            int

	Dbr          *Dbr
	Db           *db
	Duration     time.Duration
	sqlOperation SqlOperation
}

func newStmtSelect(dbr *Dbr, db *db, withStmt *StmtWith, columns *columns) *StmtSelect {
	return &StmtSelect{
		sqlOperation:      SelectOperation,
		Dbr:               dbr,
		Db:                db,
		withStmt:          withStmt,
		columns:           columns,
		conditions:        newConditions(dbr.Connections.Read),
		having:            newConditions(dbr.Connections.Read),
		distinctColumns:   newColumns(dbr.Connections.Read, false),
		distinctOnColumns: newColumns(dbr.Connections.Read, false),
		returning:         newColumns(dbr.Connections.Read, false),
	}
}

func (stmt *StmtSelect) From(tables ...interface{}) *StmtSelect {
	for _, table := range tables {
		stmt.tables = append(stmt.tables, newTable(stmt.Db, table))
	}
	return stmt
}

func (stmt *StmtSelect) Where(query interface{}, values ...interface{}) *StmtSelect {
	stmt.conditions.list = append(stmt.conditions.list, &condition{operator: OperatorAnd, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtSelect) WhereOr(query string, values ...interface{}) *StmtSelect {
	stmt.conditions.list = append(stmt.conditions.list, &condition{operator: OperatorOr, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtSelect) WhereExists(stmtSelect *StmtSelect) *StmtSelect {
	stmt.existsStmt = newStmtWhereExists(stmt.Db, stmtSelect, false)
	return stmt
}

func (stmt *StmtSelect) WhereNotExists(stmtSelect *StmtSelect) *StmtSelect {
	stmt.existsStmt = newStmtWhereExists(stmt.Db, stmtSelect, true)
	return stmt
}

func (stmt *StmtSelect) Join(table interface{}, onQuery string, values ...interface{}) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, constFunctionJoin, newTable(stmt.Db, table),
		&condition{
			operator: OperatorAnd,
			query:    onQuery,
			values:   values,
			db:       stmt.Db,
		}))
	return stmt
}

func (stmt *StmtSelect) LeftJoin(table interface{}, onQuery string, values ...interface{}) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, constFunctionLeftJoin, newTable(stmt.Db, table),
		&condition{
			operator: OperatorAnd,
			query:    onQuery,
			values:   values,
			db:       stmt.Db,
		}))
	return stmt
}

func (stmt *StmtSelect) RightJoin(table interface{}, onQuery string, values ...interface{}) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, constFunctionRightJoin, newTable(stmt.Db, table),
		&condition{
			operator: OperatorAnd,
			query:    onQuery,
			values:   values,
			db:       stmt.Db,
		}))
	return stmt
}

func (stmt *StmtSelect) FullJoin(table interface{}, onQuery string, values ...interface{}) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, constFunctionFullJoin, newTable(stmt.Db, table),
		&condition{
			operator: OperatorAnd,
			query:    onQuery,
			values:   values,
			db:       stmt.Db,
		}))
	return stmt
}

func (stmt *StmtSelect) CrossJoin(table interface{}, onQuery string, values ...interface{}) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, constFunctionCrossJoin, newTable(stmt.Db, table),
		&condition{
			operator: OperatorAnd,
			query:    onQuery,
			values:   values,
			db:       stmt.Db,
		}))
	return stmt
}

func (stmt *StmtSelect) NaturalJoin(table interface{}) *StmtSelect {
	stmt.joins = append(stmt.joins, newStmtJoin(stmt.Db, constFunctionNaturalJoin, newTable(stmt.Db, table), nil))
	return stmt
}

func (stmt *StmtSelect) Distinct(column ...interface{}) *StmtSelect {
	stmt.isDistinct = true
	stmt.distinctColumns.list = append(stmt.distinctColumns.list, column...)
	return stmt
}

func (stmt *StmtSelect) DistinctOn(column ...interface{}) *StmtSelect {
	stmt.distinctOnColumns.list = append(stmt.distinctOnColumns.list, column...)
	return stmt
}

func (stmt *StmtSelect) Union(stmtUnion *StmtSelect) *StmtSelect {
	stmt.unions = append(stmt.unions, &union{unionType: constFunctionUnion, stmt: stmtUnion})
	return stmt
}

func (stmt *StmtSelect) UnionAll(stmtUnionAll *StmtSelect) *StmtSelect {
	stmt.unions = append(stmt.unions, &union{unionType: constFunctionUnionAll, stmt: stmtUnionAll})
	return stmt
}

func (stmt *StmtSelect) Intersect(stmtIntersect *StmtSelect) *StmtSelect {
	stmt.unions = append(stmt.unions, &union{unionType: constFunctionIntersect, stmt: stmtIntersect})
	return stmt
}

func (stmt *StmtSelect) Except(stmtExcept *StmtSelect) *StmtSelect {
	stmt.unions = append(stmt.unions, &union{unionType: constFunctionExcept, stmt: stmtExcept})
	return stmt
}

func (stmt *StmtSelect) GroupBy(columns ...string) *StmtSelect {
	stmt.groupBy = append(stmt.groupBy, columns...)
	return stmt
}

func (stmt *StmtSelect) Having(query string, values ...interface{}) *StmtSelect {
	stmt.having.list = append(stmt.having.list, &condition{operator: OperatorAnd, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtSelect) HavingOr(query string, values ...interface{}) *StmtSelect {
	stmt.having.list = append(stmt.having.list, &condition{operator: OperatorOr, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtSelect) OrderBy(column string, direction direction) *StmtSelect {
	stmt.orders = append(stmt.orders, &order{column: column, direction: direction})
	return stmt
}

func (stmt *StmtSelect) OrderAsc(columns ...string) *StmtSelect {
	for _, column := range columns {
		stmt.orders = append(stmt.orders, &order{column: column, direction: OrderAsc})
	}

	return stmt
}

func (stmt *StmtSelect) OrderDesc(columns ...string) *StmtSelect {
	for _, column := range columns {
		stmt.orders = append(stmt.orders, &order{column: column, direction: OrderDesc})
	}

	return stmt
}

func (stmt *StmtSelect) Return(column ...interface{}) *StmtSelect {
	stmt.returning.list = append(stmt.returning.list, column...)
	return stmt
}

func (stmt *StmtSelect) Limit(limit int) *StmtSelect {
	stmt.limit = limit
	return stmt
}

func (stmt *StmtSelect) Offset(offset int) *StmtSelect {
	stmt.offset = offset
	return stmt
}

func (stmt *StmtSelect) Build() (string, error) {
	var query string

	// withStmt
	if len(stmt.withStmt.withs) > 0 {
		withStmt, err := stmt.withStmt.Build()
		if err != nil {
			return "", err
		}
		query += withStmt
	}

	// distinct
	var distinct string
	if stmt.isDistinct {
		distinct = fmt.Sprintf("%s ", constFunctionDistinct)
	}

	distinctColumns, err := stmt.distinctColumns.Build()
	if err != nil {
		return "", err
	}

	// distinct on
	var distinctOn string
	if stmt.isDistinct {
		distinctOn = fmt.Sprintf("%s %s ", constFunctionDistinctOn, "(%s)")
	}

	distinctOnColumns, err := stmt.distinctOnColumns.Build()
	if err != nil {
		return "", err
	}

	// columns
	columns, err := stmt.columns.Build()
	if err != nil {
		return "", err
	}

	// tables
	tables, err := stmt.tables.Build()
	if err != nil {
		return "", err
	}

	// query
	query += fmt.Sprintf("%s %s%s%s%s%s %s %s", constFunctionSelect, distinct, distinctColumns, distinctOn, distinctOnColumns, columns, constFunctionFrom, tables)

	// joins
	if len(stmt.joins) > 0 {
		joins, err := stmt.joins.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" %s", joins)
	}

	// stmt exists
	if stmt.existsStmt != nil {
		exists, err := stmt.existsStmt.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" %s %s", constFunctionWhere, exists)

	} else {

		// conditions
		if len(stmt.conditions.list) > 0 {
			conds, err := stmt.conditions.Build()
			if err != nil {
				return "", err
			}

			query += fmt.Sprintf(" %s %s", constFunctionWhere, conds)
		}
	}

	// unions
	if len(stmt.unions) > 0 {
		unions, err := stmt.unions.Build()
		if err != nil {
			return "", err
		}

		query += unions
	}

	// group by
	if len(stmt.groupBy) > 0 {
		groupBy, err := stmt.groupBy.Build()
		if err != nil {
			return "", err
		}

		query += groupBy
	}

	// having
	if len(stmt.having.list) > 0 {
		havingConds, err := stmt.having.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" %s %s", constFunctionHaving, havingConds)
	}

	// orders
	if len(stmt.orders) > 0 {
		orders, err := stmt.orders.Build()
		if err != nil {
			return "", err
		}

		query += orders
	}

	// limit
	if stmt.limit > 0 {
		query += fmt.Sprintf(" %s %d", constFunctionLimit, stmt.limit)
	}

	// offset
	if stmt.offset > 0 {
		query += fmt.Sprintf(" %s %d", constFunctionOffset, stmt.offset)
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

func (stmt *StmtSelect) Load(object interface{}) (int, error) {

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

	if err := stmt.Dbr.eventHandler(stmt.sqlOperation, stmt.tables.toArray(), query, err, rows, nil); err != nil {
		return 0, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	var cols []interface{}
	for _, col := range columns {
		cols = append(cols, col)
	}

	return read(cols, rows, value)
}
