package dbr_sync

import (
	"github.com/joaosoft/errors"
	"github.com/joaosoft/web"
)

var (
	ErrorNotFound    = errors.New(errors.LevelError, int(web.StatusNotFound), "user not found")
	ErrorInvalidType = errors.New(errors.LevelError, int(web.StatusNotFound), "invalid type")
)
