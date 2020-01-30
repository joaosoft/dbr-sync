package services

import "database/sql"

func NewExecutorDatabase(service *CmdService) *ExecutorDatabase {
	return &ExecutorDatabase{service: service}
}

type ExecutorDatabase struct {
	*sql.DB
	*sql.Tx
	service *CmdService
}

func (e *ExecutorDatabase) Open() (err error) {
	e.DB, err = e.service.config.Db.Connect()
	return err
}

func (e *ExecutorDatabase) Close() error {
	return e.DB.Close()
}

func (e *ExecutorDatabase) Begin() (err error) {
	e.Tx, err = e.DB.Begin()
	return err
}

func (e *ExecutorDatabase) Commit() error {
	return e.Tx.Commit()
}

func (e *ExecutorDatabase) Rollback() error {
	return e.Tx.Rollback()
}

func (e *ExecutorDatabase) Execute(arg interface{}, args ...interface{}) error {
	_, err := e.Tx.Exec(arg.(string), args...)
	return err
}
