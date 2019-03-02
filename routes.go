package session

import (
	"net/http"

	"github.com/joaosoft/manager"
)

func (c *Controller) RegisterRoutes(web manager.IWeb) error {
	return web.AddRoutes(
		manager.NewRoute(http.MethodGet, "/api/v1/get-session", c.GetSessionHandler),
		manager.NewRoute(http.MethodPut, "/api/v1/refresh-session", c.RefreshSessionHandler),
	)
}
