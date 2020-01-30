package services

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// CmdServiceOption ...
type CmdServiceOption func(client *CmdService)

// Reconfigure ...
func (service *CmdService) Reconfigure(options ...CmdServiceOption) {
	for _, option := range options {
		option(service)
	}
}

// WithCmdConfiguration ...
func WithCmdConfiguration(config *MigrationConfig) CmdServiceOption {
	return func(client *CmdService) {
		client.config = config
	}
}

// WithCmdLogger ...
func WithCmdLogger(logger logger.ILogger) CmdServiceOption {
	return func(service *CmdService) {
		service.logger = logger
		service.isLogExternal = true
	}
}

// WithCmdLogLevel ...
func WithCmdLogLevel(level logger.Level) CmdServiceOption {
	return func(service *CmdService) {
		service.logger.SetLevel(level)
	}
}

// WithCmdManager ...
func WithCmdManager(mgr *manager.Manager) CmdServiceOption {
	return func(service *CmdService) {
		service.pm = mgr
	}
}
