package profile

import (
	"fmt"

	"github.com/joaosoft/dbr"

	"github.com/joaosoft/manager"
	migration "github.com/joaosoft/migration/services"
)

// AppConfig ...
type AppConfig struct {
	Profile *ProfileConfig `json:"profile"`
}

// ProfileConfig ...
type ProfileConfig struct {
	Host              string                     `json:"host"`
	Dbr               *dbr.DbrConfig             `json:"dbr"`
	Migration         *migration.MigrationConfig `json:"migration"`
	Log               struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}
