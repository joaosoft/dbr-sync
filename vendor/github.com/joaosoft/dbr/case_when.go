package dbr

import (
	"fmt"
)

type caseWhen struct {
	condition *condition
	result    interface{}
}

func newCaseWhen(condition *condition) *caseWhen {
	return &caseWhen{condition: condition}
}

func (c *caseWhen) Build(db *db) (string, error) {
	var err error
	var condition string
	var result string

	// condition
	condition, err = c.condition.Build(db)
	if err != nil {
		return "", err
	}

	// result
	switch stmt := c.result.(type) {
	case *StmtSelect:
		result, err = stmt.Build()
		if err != nil {
			return "", err
		}
		result = fmt.Sprintf("(%s)", result)
	default:
		if impl, ok := stmt.(iFunction); ok {
			result, err = impl.Build(db)
			if err != nil {
				return "", err
			}
		} else {
			result = fmt.Sprintf("%+v", stmt)
		}
	}

	return fmt.Sprintf("%s %s %s %s", constFunctionWhen, condition, constFunctionThen, result), nil
}
