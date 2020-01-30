package services

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// WebServiceOption ...
type WebServiceOption func(client *WebService)

// Reconfigure ...
func (service *WebService) Reconfigure(options ...WebServiceOption) {
	for _, option := range options {
		option(service)
	}
}

// WithWebConfiguration ...
func WithWebConfiguration(config *MigrationConfig) WebServiceOption {
	return func(client *WebService) {
		client.config = config
	}
}

// WithLogger ...
func WithWebLogger(logger logger.ILogger) WebServiceOption {
	return func(service *WebService) {
		service.logger = logger
		service.isLogExternal = true
	}
}

// WithLogLevel ...
func WithWebLogLevel(level logger.Level) WebServiceOption {
	return func(service *WebService) {
		service.logger.SetLevel(level)
	}
}

// WithManager ...
func WithWebManager(mgr *manager.Manager) WebServiceOption {
	return func(service *WebService) {
		service.pm = mgr
	}
}
