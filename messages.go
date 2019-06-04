package dbr_sync

import "github.com/joaosoft/dbr"

type Operation struct {
	Operation dbr.SqlOperation `json:"operation"`
	Query     *string          `json:"query"`
	Details   *Details         `json:"details"`
}

type Details struct {
	Table      string                 `json:"table"`
	Values     map[string]interface{} `json:"values"`
	Conditions map[string]interface{} `json:"conditions"`
}

type OperationList []*Operation
