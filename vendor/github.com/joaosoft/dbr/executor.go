package dbr

import (
	"database/sql"
)

type Executor interface{
	Exec() (sql.Result, error)
}
