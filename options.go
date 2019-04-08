package profile

import (
	logger "github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// ProfileOption ...
type ProfileOption func(profile *Profile)

// Reconfigure ...
func (profile *Profile) Reconfigure(options ...ProfileOption) {
	for _, option := range options {
		option(profile)
	}
}

// WithConfiguration ...
func WithConfiguration(config *ProfileConfig) ProfileOption {
	return func(profile *Profile) {
		profile.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) ProfileOption {
	return func(profile *Profile) {
		profile.logger = logger
		profile.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) ProfileOption {
	return func(profile *Profile) {
		profile.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) ProfileOption {
	return func(profile *Profile) {
		profile.pm = mgr
	}
}
