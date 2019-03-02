package session

import (
	logger "github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// SessionOption ...
type SessionOption func(client *Session)

// Reconfigure ...
func (session *Session) Reconfigure(options ...SessionOption) {
	for _, option := range options {
		option(session)
	}
}

// WithConfiguration ...
func WithConfiguration(config *SessionConfig) SessionOption {
	return func(session *Session) {
		session.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) SessionOption {
	return func(session *Session) {
		log = logger
		session.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) SessionOption {
	return func(session *Session) {
		log.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) SessionOption {
	return func(session *Session) {
		session.pm = mgr
	}
}
