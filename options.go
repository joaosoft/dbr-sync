package dbr_sync

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// DbrSyncOption ...
type DbrSyncOption func(dbrSync *DbrSync)

// Reconfigure ...
func (dbrSync *DbrSync) Reconfigure(options ...DbrSyncOption) {
	for _, option := range options {
		option(dbrSync)
	}
}

// WithConfiguration ...
func WithConfiguration(config *DbrSyncConfig) DbrSyncOption {
	return func(dbrSync *DbrSync) {
		dbrSync.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) DbrSyncOption {
	return func(dbrSync *DbrSync) {
		dbrSync.logger = logger
		dbrSync.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) DbrSyncOption {
	return func(dbrSync *DbrSync) {
		dbrSync.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) DbrSyncOption {
	return func(dbrSync *DbrSync) {
		dbrSync.pm = mgr
	}
}
