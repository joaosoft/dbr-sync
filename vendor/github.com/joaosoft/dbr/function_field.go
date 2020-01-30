package dbr

import (
	"fmt"
)

type functionField struct {
	name      string
	field      interface{}
	arguments []interface{}

	*functionBase
}

func newFunctionField(name string, field interface{}, arguments ...interface{}) *functionField {
	return &functionField{functionBase: newFunctionBase(false, false), name: name, field: field, arguments: arguments}
}

func (c *functionField) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.name)
}

func (c *functionField) Build(db *db) (string, error) {
	c.db = db

	field, err := handleBuild(c.functionBase, c.field)
	if err != nil {
		return "", err
	}


	var arguments string
	lenArgs := len(c.arguments)
	for i, argument := range c.arguments {
		expression, err := handleBuild(c.functionBase, argument)
		if err != nil {
			return "", err
		}

		arguments += expression

		if i < lenArgs-1 {
			arguments += ", "
		}
	}

	return fmt.Sprintf("%s %s(%s)", field, c.name, arguments), nil
}
