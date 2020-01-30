package services

import (
	"sync"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type WebService struct {
	config        *MigrationConfig
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
	logger        logger.ILogger
}

// NewWebService ...
func NewWebService(options ...WebServiceOption) (*WebService, error) {
	config, simpleConfig, err := NewConfig()
	service := &WebService{
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		logger: logger.NewLogDefault("migration", logger.WarnLevel),
		config: config.Migration,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Migration != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Migration.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	} else {
		config.Migration = &MigrationConfig{
			Host: DefaultURL,
		}
	}

	service.Reconfigure(options...)

	simpleDB := service.pm.NewSimpleDB(&config.Migration.Db.DBConfig)
	if err := service.pm.AddDB("db_postgres", simpleDB); err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	web := service.pm.NewSimpleWebEcho(service.config.Host)
	controller := NewController(service.logger, NewInteractor(service.logger, NewStoragePostgres(service.logger, simpleDB)))
	controller.RegisterRoutes(web)

	service.pm.AddWeb("api_web", web)

	return service, nil
}

// Start ...
func (m *WebService) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *WebService) Stop() error {
	return m.pm.Stop()
}
