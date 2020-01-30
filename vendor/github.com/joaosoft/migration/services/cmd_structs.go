package services

type Tag string
type MigrationCommand string
type MigrationOption string
type CustomMode string

type Handler func(option MigrationOption, conn Executor, data string) error
