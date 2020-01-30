package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type CmdService struct {
	config        *MigrationConfig
	interactor    *Interactor
	tag           map[string]Handler
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
	logger        logger.ILogger
}

func NewCmdService(options ...CmdServiceOption) (*CmdService, error) {
	config, simpleConfig, err := NewConfig()

	service := &CmdService{
		pm:     manager.NewManager(manager.WithRunInBackground(true)),
		logger: logger.NewLogDefault("migration", logger.WarnLevel),
		tag: map[string]Handler{
			string(FileTagMigrateUp):   MigrationHandler,
			string(FileTagMigrateDown): MigrationHandler,
		},
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

	simpleDB := service.pm.NewSimpleDB(&service.config.Db.DBConfig)
	if err := service.pm.AddDB("db_postgres", simpleDB); err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	simpleDB.Start(wg)
	service.interactor = NewInteractor(service.logger, NewStoragePostgres(service.logger, simpleDB))

	return service, nil
}

func (service *CmdService) AddTag(name string, handler Handler) error {
	_, okUp := service.tag[fmt.Sprintf(string(FileTagCustomUp), name)]
	_, okDown := service.tag[fmt.Sprintf(string(FileTagCustomDown), name)]

	if okUp || okDown {
		return service.logger.Errorf("the tag %s already exists!", name).ToError()
	}

	service.tag[fmt.Sprintf(string(FileTagCustomUp), name)] = handler
	service.tag[fmt.Sprintf(string(FileTagCustomDown), name)] = handler

	return nil
}

// execute ...
func (service *CmdService) Execute(option MigrationOption, number int, mode ExecutorMode) (int, error) {
	service.logger.Infof("executing migration with option '-%s %s'", CmdMigrate, option)

	// setup
	if err := service.setup(); err != nil {
		return 0, err
	}

	// load
	executed, toexecute, err := service.load(mode)
	if err != nil {
		return 0, err
	}

	// validate
	if err := service.validate(executed, toexecute); err != nil {
		return 0, err
	}

	// process
	return service.process(option, number, mode, executed, toexecute)
}

// setup ...
func (service *CmdService) setup() error {
	service.logger.Info("executing setup of migration table")

	conn, err := service.config.Db.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(fmt.Sprintf(CREATE_MIGRATION_TABLES, service.config.Db.Schema)); err != nil {
		service.logger.Error("error creating migration table")
	}

	if err = tx.Commit(); err != nil {
		service.logger.Error("error executing commit of transaction")
		return err
	}

	return err
}

// load ...
func (service *CmdService) load(mode ExecutorMode) (executed []string, toexecute []string, err error) {

	// executed
	migrations, er := service.interactor.GetMigrations(map[string][]string{"mode": []string{string(mode)}})
	if er != nil {
		return nil, nil, service.logger.Error("error loading migrations from database").ToError()
	}
	for _, migration := range migrations {
		executed = append(executed, migration.IdMigration)
	}

	// to execute
	dir, _ := os.Getwd()

	files, err := filepath.Glob(fmt.Sprintf("%s/%s/*", dir, service.getPath(mode)))
	if err != nil {
		return executed, nil, service.logger.Error("error loading migrations from file system").ToError()
	}
	for _, file := range files {
		fileName := file[strings.LastIndex(file, "/")+1:]
		toexecute = append(toexecute, fileName)
	}

	return executed, toexecute, err
}

// validate ...
func (service *CmdService) validate(executed []string, toexecute []string) (err error) {

	if len(executed) > len(toexecute) {
		return service.logger.Errorf("error, the migration number is different [executed: %d] <-> [to execute: %d]", len(executed), len(toexecute)).ToError()
	}

	for i, migration := range executed {
		if migration != toexecute[i] {
			return service.logger.Errorf("error, the migrations are in a different order of the already executed migrations [%s] <-> [%s]", migration, toexecute[i]).ToError()
		}
	}
	return nil
}

// process ...
func (service *CmdService) process(option MigrationOption, number int, mode ExecutorMode, executed []string, toexecute []string) (int, error) {
	var migrations []string

	if option == OptionUp {
		if len(toexecute) <= len(executed) {
			service.logger.Infof("applied %d migrations!", 0)
			return 0, nil
		}

		if number > (len(toexecute) - len(executed)) {
			number = len(toexecute) - len(executed)
		}
		sort.Slice(toexecute, func(i, j int) bool {
			if toexecute[i] < toexecute[j] {
				return true
			}
			return false
		})

		if number > 0 {
			migrations = toexecute[len(executed) : len(executed)+number]
		} else {
			migrations = toexecute[len(executed):]
		}
	} else {
		if len(executed) == 0 {
			service.logger.Infof("applied %d migrations!", 0)
			return 0, nil
		}
		toexecute = toexecute[:len(executed)]
		sort.Slice(toexecute, func(i, j int) bool {
			if toexecute[i] < toexecute[j] {
				return false
			}
			return true
		})

		if number == 0 {
			number = 1
		}

		if number > 0 {
			migrations = toexecute[:number]
		} else {
			migrations = toexecute
		}
	}

	executor := NewExecutor(service, mode)
	if err := executor.Open(); err != nil {
		service.logger.Error("error connecting executor")
		return 0, err
	}

	defer executor.Close()

	// log migration information
	service.logger.Infof("migrations already executed %+v", executed)
	service.logger.Infof("migrations to execute %+v", migrations)

	for _, migration := range migrations {
		migrationTags, customTags, err := service.loadRunningTags(option, mode, migration)
		if err != nil {
			return 0, err
		}

		if err = executor.Begin(); err != nil {
			return 0, err
		}

		// execute migration handlers
		for key, value := range migrationTags {
			if err = service.tag[key](option, executor, value); err != nil {
				break
			}
		}

		if err == nil {
			// execute custom handlers
			for key, value := range customTags {
				if err = service.tag[key](option, executor, value); err != nil {
					break
				}
			}
		}

		if err == nil {
			if err = executor.Commit(); err != nil {
				service.logger.Error("error executing commit of transaction")
				return 0, err
			}
		}

		switch option {
		case OptionUp:
			if err == nil {
				if err = service.interactor.CreateMigration(&Migration{IdMigration: migration, Mode: mode}); err != nil {
					service.logger.Error("error adding migration to database")
					service.interactor.DeleteMigration(migration)
					executor.Rollback()
					return 0, err
				}
			}
		default:
			if err == nil {
				if err = service.interactor.DeleteMigration(migration); err != nil {
					service.logger.Error("error deleting migration to database")
					executor.Rollback()
					return 0, err
				}
			}
		}

		if err != nil {
			service.logger.Errorf("error executing the migration %s", migration)
			executor.Rollback()
			return 0, err
		}
	}
	service.logger.Infof("applied %d migrations!", len(migrations))

	return len(migrations), nil
}

func (service *CmdService) loadRunningTags(option MigrationOption, mode ExecutorMode, file string) (migrationTags map[string]string, customTags map[string]string, err error) {
	dir, _ := os.Getwd()
	lines, err := ReadFileLines(fmt.Sprintf("%s/%s/%s", dir, service.getPath(mode), file))
	if err != nil {
		return nil, nil, err
	}

	var tag string
	var text string

	migrationTags = make(map[string]string)
	customTags = make(map[string]string)

	addFunc := func(tag string, text *string, migrationTags, customTags map[string]string) {
		if tag != "" && *text != "" {
			if tag == fmt.Sprintf(string(FileTagMigrate), option) {
				migrationTags[tag] = *text
			} else {
				if strings.HasSuffix(tag, string(option)) {
					customTags[tag] = *text
				}
			}
			*text = ""
		}
	}

	for _, line := range lines {
		if _, ok := service.tag[strings.TrimSpace(line)]; ok {
			addFunc(tag, &text, migrationTags, customTags)
			tag = strings.TrimSpace(line)
			continue
		}
		text += fmt.Sprintf("%s\n", line)
	}

	addFunc(tag, &text, migrationTags, customTags)

	return migrationTags, customTags, nil
}

// Start ...
func (m *CmdService) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *CmdService) Stop() error {
	return m.pm.Stop()
}

func (m *CmdService) getPath(mode ExecutorMode) string {
	switch mode {
	case ExecutorModeDatabase:
		return m.config.Path.Database
	case ExecutorModeRabbitMq:
		return m.config.Path.Rabbitmq
	}

	return ""
}
