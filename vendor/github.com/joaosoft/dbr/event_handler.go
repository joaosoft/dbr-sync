package dbr

import "database/sql"

type eventHandler func(operation SqlOperation, table []string, query string, err error, rows *sql.Rows, sqlResult sql.Result) error

type SuccessEventHandler func(operation SqlOperation, table []string, query string, rows *sql.Rows, sqlResult sql.Result) error
type ErrorEventHandler func(operation SqlOperation, table []string, query string, err error) error

func NewDb(database database, dialect dialect) *db {
	return &db{
		database: database,
		Dialect:  dialect,
	}
}

func (dbr *Dbr) handle(operation SqlOperation, table []string, query string, err error, rows *sql.Rows, sqlResult sql.Result) error {
	if err == nil && dbr.successEventHandler != nil {
		if err := dbr.successEventHandler(operation, table, query, rows, sqlResult); err != nil {
			return err
		}
	}

	if err != nil && dbr.errorEventHandler != nil {
		if err := dbr.errorEventHandler(operation, table, query, err); err != nil {
			return err
		}
	}

	return nil
}
