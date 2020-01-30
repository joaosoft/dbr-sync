package services

type ExecutorMode string

const (
	ExecutorModeAll      ExecutorMode = "all"
	ExecutorModeDatabase ExecutorMode = "database"
	ExecutorModeRabbitMq ExecutorMode = "rabbitmq"
)

type Executor interface {
	Open() error
	Begin() error
	Execute(arg interface{}, args ...interface{}) error
	Commit() error
	Rollback() error
	Close() error
}

func NewExecutor(service *CmdService, mode ExecutorMode) Executor {
	switch mode {
	case ExecutorModeDatabase:
		return NewExecutorDatabase(service)
	case ExecutorModeRabbitMq:
		return NewExecutorRabbitMq(service)
	}

	return nil
}
