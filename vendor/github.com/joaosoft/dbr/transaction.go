package dbr

import (
	"database/sql"
	"time"
)

type Transaction struct {
	commited  bool
	dbr       *Dbr
	db        *db
	startTime time.Time
	Duration  time.Duration
}

func newTransaction(dbr *Dbr, db *db, startTime time.Time) *Transaction {
	return &Transaction{dbr: dbr, db: db, startTime: startTime}
}

func (tx *Transaction) Commit() error {

	defer func() {
		tx.Duration = time.Since(tx.startTime)
	}()

	err := tx.db.database.(*sql.Tx).Commit()
	if err != nil {
		return err
	}

	tx.commited = true
	return nil
}

func (tx *Transaction) Rollback() error {

	defer func() {
		tx.Duration = time.Since(tx.startTime)
	}()

	return tx.db.database.(*sql.Tx).Rollback()
}

func (tx *Transaction) RollbackUnlessCommit() error {
	if !tx.commited {
		return tx.db.database.(*sql.Tx).Rollback()
	}

	return nil
}

func (tx *Transaction) Select(column ...interface{}) *StmtSelect {
	columns := newColumns(tx.db, false)
	columns.list = column

	return newStmtSelect(tx.dbr, tx.dbr.Connections.Write, &StmtWith{}, columns)
}

func (tx *Transaction) Insert() *StmtInsert {
	return newStmtInsert(tx.dbr, tx.dbr.Connections.Write, &StmtWith{})
}

func (tx *Transaction) Update(table string) *StmtUpdate {
	return newStmtUpdate(tx.dbr, tx.dbr.Connections.Write, &StmtWith{}, table)
}

func (tx *Transaction) Delete() *StmtDelete {
	return newStmtDelete(tx.dbr, tx.dbr.Connections.Write, &StmtWith{})
}

func (tx *Transaction) Execute(query string) *StmtExecute {
	return newStmtExecute(tx.dbr, tx.dbr.Connections.Write, query)
}

func (tx *Transaction) With(name string, builder Builder) *StmtWith {
	return newStmtWith(tx.dbr, tx.dbr.Connections, name, false, builder)
}

func (tx *Transaction) WithRecursive(name string, builder Builder) *StmtWith {
	return newStmtWith(tx.dbr, tx.dbr.Connections, name, true, builder)
}
