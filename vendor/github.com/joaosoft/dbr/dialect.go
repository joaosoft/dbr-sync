package dbr

import "time"

type dialectName string

var (
	availableDialects = map[dialectName]dialect {
		constDialectPostgres: &dialectPostgres{},
		constDialectMysql: &dialectMySql{},
		constDialectSqlLite3: &dialectSqlLite3{},
	}
)

type dialect interface {
	Name() string
	Encode(value interface{}) string
	EncodeString(value string) string
	EncodeBool(value bool) string
	EncodeTime(value time.Time) string
	EncodeBytes(value []byte) string
	EncodeColumn(column interface{}) string
	Placeholder() string
}

func newDialect(name dialectName) (dialect, error) {
	dialect, found := availableDialects[name]
	if !found {
		return nil, ErrorDialectNotFound
	}

	return dialect, nil
}
