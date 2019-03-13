package session

import (
	"sync"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	migration "github.com/joaosoft/migration/services"
)

type Session struct {
	config        *SessionConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
	mux           sync.Mutex
}

// NewSession ...
func NewSession(options ...SessionOption) (*Session, error) {
	config, simpleConfig, err := NewConfig()

	service := &Session{
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		logger: logger.NewLogDefault("session", logger.WarnLevel),
		config: config.Session,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(logger.Instance))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Session != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Session.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	} else {
		config.Session = &SessionConfig{
			Host: defaultURL,
		}
	}

	service.Reconfigure(options...)

	// execute migrations
	migrationService, err := migration.NewCmdService(migration.WithCmdConfiguration(service.config.Migration))
	if err != nil {
		return nil, err
	}

	if _, err := migrationService.Execute(migration.OptionUp, 0, migration.ExecutorModeDatabase); err != nil {
		return nil, err
	}

	web := service.pm.NewSimpleWebServer(config.Session.Host)

	storage, err := NewStoragePostgres(config.Session)
	if err != nil {
		return nil, err
	}

	interactor := NewInteractor(config.Session, storage)

	controller := NewController(config.Session, interactor)
	controller.RegisterRoutes(web)

	service.pm.AddWeb("api_web", web)

	return service, nil
}

// Start ...
func (m *Session) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *Session) Stop() error {
	return m.pm.Stop()
}
