package services

func MigrationHandler(option MigrationOption, conn Executor, data string) error {
	err := conn.Execute(data)

	return err
}
