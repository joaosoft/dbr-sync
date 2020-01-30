package dbr

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	"github.com/joaosoft/migration/services"
)

// DbrOption ...
type DbrOption func(client *Dbr)

// Reconfigure ...
func (dbr *Dbr) Reconfigure(options ...DbrOption) {
	for _, option := range options {
		option(dbr)
	}
}

// WithConfiguration ...
func WithConfiguration(config *DbrConfig) DbrOption {
	return func(dbr *Dbr) {
		dbr.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) DbrOption {
	return func(dbr *Dbr) {
		dbr.logger = logger
		dbr.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) DbrOption {
	return func(dbr *Dbr) {
		dbr.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) DbrOption {
	return func(dbr *Dbr) {
		dbr.pm = mgr
	}
}

// WithDatabase ...
func WithDatabase(master *db, slave ...*db) DbrOption {
	return func(dbr *Dbr) {
		if len(slave) > 0 {
			dbr.Connections = &connections{Write: master, Read: slave[0]}
		} else {
			dbr.Connections = &connections{Write: master, Read: master}
		}
	}
}

// WithMigrationConfig ...
func WithMigrationConfig(migration *services.MigrationConfig) DbrOption {
	return func(dbr *Dbr) {
		dbr.config.Migration = migration
	}
}

// WithSuccessEventHandler ...
func WithSuccessEventHandler(eventHandler SuccessEventHandler) DbrOption {
	return func(dbr *Dbr) {
		dbr.successEventHandler = eventHandler
	}
}

// WithErrorEventHandler ...
func WithErrorEventHandler(eventHandler ErrorEventHandler) DbrOption {
	return func(dbr *Dbr) {
		dbr.errorEventHandler = eventHandler
	}
}
