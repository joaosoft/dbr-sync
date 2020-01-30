package services

import (
	"fmt"

	"github.com/joaosoft/errors"

	"github.com/joaosoft/web"
)

func NewExecutorRabbitMq(service *CmdService) *ExecutorRabbitMq {
	return &ExecutorRabbitMq{service: service}
}

type ExecutorRabbitMq struct {
	client  *web.Client
	service *CmdService
}

func (e *ExecutorRabbitMq) Open() (err error) {
	e.client, err = web.NewClient()
	return err
}

func (e *ExecutorRabbitMq) Close() error {
	return nil
}

func (e *ExecutorRabbitMq) Begin() error {
	return nil
}

func (e *ExecutorRabbitMq) Commit() error {
	return nil
}

func (e *ExecutorRabbitMq) Rollback() error {
	return nil
}

func (e *ExecutorRabbitMq) Execute(arg interface{}, args ...interface{}) error {
	url := fmt.Sprintf("%s/api/definitions", e.service.config.RabbitMq.Host)

	if e.service.config.RabbitMq.VHost != nil {
		url += fmt.Sprintf("/%s", *e.service.config.RabbitMq.VHost)
	}

	request, err := e.client.NewRequest(web.MethodPost, url, web.ContentTypeApplicationJSON, nil)
	if err != nil {
		return err
	}

	request.Headers["Authorization"] = []string{"Basic Zm91cnNvdXJjZTpmNHMwdTQ1ZQ=="}

	response, err := request.WithBody([]byte(arg.(string))).Send()
	if err != nil {
		return err
	}

	if response.Status >= web.StatusBadRequest {
		return errors.New(errors.ErrorLevel, 0, "error importing configurations to rabbitmq [status: %d, error: %s]", response.Status, string(response.Body))
	}

	return nil
}
