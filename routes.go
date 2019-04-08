package profile

import (
	"github.com/joaosoft/manager"
	"github.com/joaosoft/web"
)

func (c *Controller) RegisterRoutes(w manager.IWeb) error {
	return w.AddRoutes(
		manager.NewRoute(string(web.MethodGet), "/api/v1/profile/sections", c.GetSectionsHandler),
		manager.NewRoute(string(web.MethodGet), "/api/v1/profile/sections/:section_key", c.GetSectionHandler),
		manager.NewRoute(string(web.MethodGet), "/api/v1/profile/sections/:section_key/contents", c.GetSectionContentsHandler),
	)
}
