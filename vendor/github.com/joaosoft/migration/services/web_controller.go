package services

import (
	"net/http"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/validator"
	"github.com/labstack/echo"
)

type Controller struct {
	logger     logger.ILogger
	interactor *Interactor
}

func NewController(logger logger.ILogger, interactor *Interactor) *Controller {
	return &Controller{
		logger:     logger,
		interactor: interactor,
	}
}

func (controller *Controller) GetMigrationHandler(ctx echo.Context) error {
	request := GetMigrationRequest{
		IdMigration: ctx.Param("id"),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		return ctx.JSON(http.StatusBadRequest, errs)
	}

	if process, err := controller.interactor.GetMigration(request.IdMigration); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if process == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, process)
	}
}

func (controller *Controller) GetMigrationsHandler(ctx echo.Context) error {
	if processes, err := controller.interactor.GetMigrations(ctx.QueryParams()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if processes == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, processes)
	}
}

func (controller *Controller) CreateMigrationHandler(ctx echo.Context) error {
	request := CreateMigrationRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		err = controller.logger.WithFields(map[string]interface{}{"error": err}).
			Error("error getting body").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if errs := validator.Validate(request.Body); len(errs) > 0 {
		newErr := errors.New(errors.ErrorLevel, 0, errs)
		controller.logger.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	newMigration := Migration{
		IdMigration: request.Body.IdMigration,
	}
	if err := controller.interactor.CreateMigration(&newMigration); err != nil {
		newErr := errors.New(errors.ErrorLevel, 0, err)
		controller.logger.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Errorf("error creating process %s", request.Body.IdMigration).ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		return ctx.NoContent(http.StatusCreated)
	}
}

func (controller *Controller) DeleteMigrationHandler(ctx echo.Context) error {
	request := DeleteMigrationRequest{
		IdMigration: ctx.Param("id"),
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		newErr := errors.New(errors.ErrorLevel, 0, errs)
		controller.logger.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := controller.interactor.DeleteMigration(request.IdMigration); err != nil {
		newErr := errors.New(errors.ErrorLevel, 0, err)
		controller.logger.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Errorf("error deleting process by id %s", request.IdMigration).ToError()
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (controller *Controller) DeleteMigrationsHandler(ctx echo.Context) error {
	if err := controller.interactor.DeleteMigrations(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}
