package dbr_sync

import (
	"fmt"

	"github.com/joaosoft/dbr"

	"github.com/joaosoft/manager"
	migration "github.com/joaosoft/migration/services"
)

// AppConfig ...
type AppConfig struct {
	DbrSync *DbrSyncConfig `json:"dbr-sync"`
}

type RabbitmqConfig struct {
	*manager.RabbitmqConfig
	Queue   string `json:"queue"`
	Binding string `json:"binding"`
}

// ProfileConfig ...
type DbrSyncConfig struct {
	Rabbitmq  RabbitmqConfig             `json:"rabbitmq"`
	Dbr       *dbr.DbrConfig             `json:"dbr"`
	Migration *migration.MigrationConfig `json:"migration"`
	Mode      mode                       `json:"mode"`
	Log       struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}
