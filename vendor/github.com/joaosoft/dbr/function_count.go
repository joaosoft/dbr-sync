package dbr

import (
	"fmt"
)

type functionCount struct {
	expression interface{}
	distinct   bool

	*functionBase
}

func newFunctionCount(expression interface{}, distinct bool) *functionCount {
	return &functionCount{functionBase: newFunctionBase(false, false), expression: expression, distinct: distinct}
}

func (c *functionCount) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionCount) Build(db *db) (string, error) {
	c.db = db

	expression, err := handleBuild(c.functionBase, c.expression)
	if err != nil {
		return "", err
	}

	var distinct string
	if c.distinct {
		distinct = fmt.Sprintf("%s ", constFunctionDistinct)
	}

	return fmt.Sprintf("%s(%s%s)", constFunctionCount, distinct, expression), nil
}
