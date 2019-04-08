package profile

import (
	"sync"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	migration "github.com/joaosoft/migration/services"
)

type Profile struct {
	config        *ProfileConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
	mux           sync.Mutex
}

// NewProfile ...
func NewProfile(options ...ProfileOption) (*Profile, error) {
	config, simpleConfig, err := NewConfig()

	service := &Profile{
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		logger: logger.NewLogDefault("profile", logger.WarnLevel),
		config: config.Profile,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(logger.Instance))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Profile != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Profile.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	} else {
		config.Profile = &ProfileConfig{
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

	web := service.pm.NewSimpleWebServer(config.Profile.Host)

	storage, err := NewStoragePostgres(config.Profile)
	if err != nil {
		return nil, err
	}

	interactor := NewInteractor(service.logger, storage)

	controller := NewController(config.Profile, interactor)
	controller.RegisterRoutes(web)

	service.pm.AddWeb("api_web", web)

	return service, nil
}

// Start ...
func (m *Profile) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *Profile) Stop() error {
	return m.pm.Stop()
}
