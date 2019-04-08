package profile

import (
	"github.com/joaosoft/validator"
	"github.com/joaosoft/web"
)

type Controller struct {
	config     *ProfileConfig
	interactor *Interactor
}

func NewController(config *ProfileConfig, interactor *Interactor) *Controller {
	return &Controller{
		config:     config,
		interactor: interactor,
	}
}

func (c *Controller) GetSectionsHandler(ctx *web.Context) error {
	response, err := c.interactor.GetSections()
	if err != nil {
		return ctx.Response.JSON(web.StatusInternalServerError, ErrorResponse{Code: web.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Response.JSON(web.StatusOK, response)
}

func (c *Controller) GetSectionHandler(ctx *web.Context) error {
	request := &GetSectionRequest{}

	if err := ctx.Request.BindUrlParams(&request); err != nil {
		return ctx.Response.JSON(web.StatusBadRequest, err)
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		return ctx.Response.JSON(web.StatusBadRequest, errs)
	}

	response, err := c.interactor.GetSection(request)
	if err != nil {
		return ctx.Response.JSON(web.StatusInternalServerError, ErrorResponse{Code: web.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Response.JSON(web.StatusOK, response)
}

func (c *Controller) GetSectionContentsHandler(ctx *web.Context) error {
	request := &GetSectionContentsRequest{}

	if err := ctx.Request.BindUrlParams(&request); err != nil {
		return ctx.Response.JSON(web.StatusBadRequest, err)
	}

	if errs := validator.Validate(request); len(errs) > 0 {
		return ctx.Response.JSON(web.StatusBadRequest, errs)
	}

	response, err := c.interactor.GetSectionContents(request)
	if err != nil {
		return ctx.Response.JSON(web.StatusInternalServerError, ErrorResponse{Code: web.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.Response.JSON(web.StatusOK, response)
}
