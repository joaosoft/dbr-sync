package dbr_sync

import (
	"fmt"
	"sync"

	"github.com/joaosoft/dbr"
	"github.com/joaosoft/json"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	migration "github.com/joaosoft/migration/services"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
)

type DbrSync struct {
	config        *DbrSyncConfig
	isLogExternal bool
	pm            *manager.Manager
	storage       *dbr.Dbr
	logger        logger.ILogger
	mux           sync.Mutex
}

// NewDbrSync ...
func NewDbrSync(options ...DbrSyncOption) (*DbrSync, error) {
	config, simpleConfig, err := NewConfig()

	service := &DbrSync{
		pm:     manager.NewManager(manager.WithRunInBackground(true)),
		logger: logger.NewLogDefault("dbr-sync", logger.WarnLevel),
		config: config.DbrSync,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(logger.Instance))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.DbrSync != nil {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.DbrSync.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	} else {
		config.DbrSync = &DbrSyncConfig{}
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

	service.storage, err = dbr.New(dbr.WithConfiguration(config.DbrSync.Dbr))
	if err != nil {
		return nil, err
	}

	// rabbitmq
	rabbitmqConsumer, err := service.pm.NewSimpleRabbitmqConsumer(
		config.DbrSync.Rabbitmq.RabbitmqConfig,
		config.DbrSync.Rabbitmq.Queue,
		config.DbrSync.Rabbitmq.Binding,
		"dbr-sync", service.consumeRabbitMessage)
	if err != nil {
		log.Errorf("%s", err)
	}

	service.pm.AddRabbitmqConsumer("rabbitmq_consumer", rabbitmqConsumer)

	return service, nil
}

// Start ...
func (m *DbrSync) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *DbrSync) Stop() error {
	return m.pm.Stop()
}

func (m *DbrSync) consumeRabbitMessage(msg amqp.Delivery) error {
	log.Info("Handling message...")
	var operations OperationList

	err := json.Unmarshal(msg.Body, &operations, "json")
	if err != nil {
		return err
	}

	var connection dbr.IDbr
	connection = m.storage
	var executor dbr.Executor

	lenOperations := len(operations)

	if lenOperations > 1 {
		tx, err := m.storage.Begin()
		if err != nil {
			return err
		}

		connection = tx
	}

	for _, operation := range operations {

		if operation.Query != nil {
			executor = connection.
				Execute(*operation.Query)
		} else {
			switch operation.Operation {
			case dbr.InsertOperation:

				// columns / values
				var columns []interface{}
				var values []interface{}
				for column, value := range operation.Details.Values {
					columns = append(columns, column)
					values = append(values, value)
				}

				executor = connection.
					Insert().
					Into(operation.Details.Table).
					Columns(columns...).
					Values(values...)

			case dbr.UpdateOperation:
				stmtUpdate := connection.
					Update(operation.Details.Table)

				for column, value := range operation.Details.Values {
					stmtUpdate.Set(column, value)
				}

				// query
				var query string
				var values []interface{}
				lenC := len(operation.Details.Conditions)

				for column, value := range operation.Details.Conditions {
					query += fmt.Sprintf("%s = ?", column)
					values = append(values, value)

					if len(values) < lenC {
						query += ", "
					}
				}

				stmtUpdate.Where(query, values...)

				executor = stmtUpdate

			case dbr.DeleteOperation:
				stmtDelete := connection.
					Delete().From(operation.Details.Table)

				// query
				var query string
				var values []interface{}
				lenC := len(operation.Details.Conditions)

				for column, value := range operation.Details.Conditions {
					query += fmt.Sprintf("%s = ?", column)
					values = append(values, value)

					if len(values) < lenC {
						query += ", "
					}
				}

				stmtDelete.Where(query, values...)

				executor = stmtDelete
			}
		}

		_, err = executor.Exec()
		if err != nil {
			return err
		}
	}

	if lenOperations > 1 {
		return connection.(*dbr.Transaction).Commit()
	}

	return nil
}
