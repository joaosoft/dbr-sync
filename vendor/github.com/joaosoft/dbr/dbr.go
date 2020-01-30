package dbr

import (
	"database/sql"
	"sync"
	"time"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	"github.com/joaosoft/migration/services"
)

type Dbr struct {
	Connections *connections

	eventHandler        eventHandler
	successEventHandler SuccessEventHandler
	errorEventHandler   ErrorEventHandler
	config              *DbrConfig
	logger              logger.ILogger
	isLogExternal       bool
	pm                  *manager.Manager
	mux                 sync.Mutex
}

type IDbr interface {
	Select(column ...interface{}) *StmtSelect
	Insert() *StmtInsert
	Update(table string) *StmtUpdate
	Delete() *StmtDelete
	Execute(query string) *StmtExecute
	With(name string, builder Builder) *StmtWith
	WithRecursive(name string, builder Builder) *StmtWith
}

type connections struct {
	Read  *db
	Write *db
}

type db struct {
	database
	Dialect dialect
}

// New ...
func New(options ...DbrOption) (*Dbr, error) {
	config, simpleConfig, err := NewConfig()

	service := &Dbr{
		pm:     manager.NewManager(manager.WithRunInBackground(true)),
		logger: logger.NewLogDefault("dbr", logger.WarnLevel),
		config: config.Dbr,
	}

	// set the internal event handler
	service.eventHandler = service.handle

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Dbr != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Dbr.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	service.Reconfigure(options...)

	// connect to database
	if service.config != nil && service.Connections == nil {
		if service.config.Db != nil {
			dbCon := service.pm.NewSimpleDB(service.config.Db)
			if err := dbCon.Start(); err != nil {
				return nil, err
			}
			service.pm.AddDB("db", dbCon)

			dbDialect, err := newDialect(dialectName(service.config.Db.Driver))
			if err != nil {
				return nil, err
			}
			db := &db{database: dbCon.Get(), Dialect: dbDialect}

			service.Connections = &connections{Read: db, Write: db}
		} else if service.config.ReadDb != nil && service.config.WriteDb != nil {
			dbReadCon := service.pm.NewSimpleDB(service.config.ReadDb)
			if err := dbReadCon.Start(); err != nil {
				return nil, err
			}
			service.pm.AddDB("db-read", dbReadCon)

			dbReadDialect, err := newDialect(dialectName(service.config.ReadDb.Driver))
			if err != nil {
				return nil, err
			}
			dbRead := &db{database: dbReadCon.Get(), Dialect: dbReadDialect}

			dbWriteCon := service.pm.NewSimpleDB(service.config.WriteDb)
			if err := dbWriteCon.Start(); err != nil {
				return nil, err
			}
			service.pm.AddDB("db-write", dbWriteCon)

			dbWriteDialect, err := newDialect(dialectName(service.config.WriteDb.Driver))
			if err != nil {
				return nil, err
			}
			dbWrite := &db{database: dbReadCon.Get(), Dialect: dbWriteDialect}

			service.Connections = &connections{Read: dbRead, Write: dbWrite}
		}

		// execute migrations
		if service.config.Migration != nil {
			migration, err := services.NewCmdService(services.WithCmdConfiguration(service.config.Migration))
			if err != nil {
				return nil, err
			}

			if _, err := migration.Execute(services.OptionUp, 0, services.ExecutorModeDatabase); err != nil {
				return nil, err
			}
		}
	}

	return service, nil
}

func (dbr *Dbr) Select(column ...interface{}) *StmtSelect {
	columns := newColumns(dbr.Connections.Read, false)
	columns.list = column

	return newStmtSelect(dbr, dbr.Connections.Read, &StmtWith{}, columns)
}

func (dbr *Dbr) Insert() *StmtInsert {
	return newStmtInsert(dbr, dbr.Connections.Write, &StmtWith{})
}

func (dbr *Dbr) Update(table string) *StmtUpdate {
	return newStmtUpdate(dbr, dbr.Connections.Write, &StmtWith{}, table)
}

func (dbr *Dbr) Delete() *StmtDelete {
	return newStmtDelete(dbr, dbr.Connections.Write, &StmtWith{})
}

func (dbr *Dbr) Execute(query string) *StmtExecute {
	return newStmtExecute(dbr, dbr.Connections.Write, query)
}

func (dbr *Dbr) With(name string, builder Builder) *StmtWith {
	return newStmtWith(dbr, dbr.Connections, name, false, builder)
}

func (dbr *Dbr) WithRecursive(name string, builder Builder) *StmtWith {
	return newStmtWith(dbr, dbr.Connections, name, true, builder)
}

func (dbr *Dbr) UseOnlyWrite(name string, builder Builder) *Dbr {
	return &Dbr{
		config:        dbr.config,
		logger:        dbr.logger,
		isLogExternal: dbr.isLogExternal,
		pm:            dbr.pm,
		mux:           dbr.mux,
		Connections: &connections{
			Read:  dbr.Connections.Write,
			Write: dbr.Connections.Write,
		},
	}
}

func (dbr *Dbr) UseOnlyRead(name string, builder Builder) *Dbr {
	return &Dbr{
		config:        dbr.config,
		logger:        dbr.logger,
		isLogExternal: dbr.isLogExternal,
		pm:            dbr.pm,
		mux:           dbr.mux,
		Connections:   &connections{Read: dbr.Connections.Read, Write: dbr.Connections.Read},
	}
}

func (dbr *Dbr) Begin() (*Transaction, error) {
	startTime := time.Now()
	tx, err := dbr.Connections.Write.database.(*sql.DB).Begin()
	if err != nil {
		return nil, err
	}

	return newTransaction(dbr, &db{database: tx, Dialect: dbr.Connections.Write.Dialect}, startTime), nil
}
