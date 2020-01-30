package dbr

import (
	"fmt"
	"strings"
)

type functionCase struct {
	alias   *string
	onCase  *onCase
	onWhens onCaseWhens
	onElse  *onCaseElse

	*functionBase
}

func newFunctionCase(value ...interface{}) *functionCase {
	funcCase := &functionCase{
		functionBase: newFunctionBase(false, false),
		onCase:       newCase(value...),
		onWhens:      newCaseWhens(),
	}

	return funcCase
}

func (c *functionCase) When(query interface{}, values ...interface{}) *functionCase {
	c.onWhens = append(c.onWhens, newCaseWhen(newCondition(nil, OperatorAnd, query, values...)))

	return c
}

func (c *functionCase) Then(result interface{}) *functionCase {
	if len(c.onWhens) > 0 {
		c.onWhens[len(c.onWhens)-1].result = result
	}

	return c
}

func (c *functionCase) Else(result interface{}) *functionCase {
	c.onElse = newCaseElse(result)

	return c
}

func (c *functionCase) As(alias string) *functionCase {
	c.alias = &alias

	return c
}

func (c *functionCase) Expression(db *db) (string, error) {
	c.db = db
	return "", nil
}

func (c *functionCase) Build(db *db) (string, error) {
	c.db = db

	var value string
	var query string

	onCase, err := c.onCase.Build(db)
	if err != nil {
		return "", err
	}
	value += onCase

	onWhens, err := c.onWhens.Build(db)
	if err != nil {
		return "", err
	}
	value += onWhens

	onElse, err := c.onElse.Build(db)
	if err != nil {
		return "", err
	}

	if len(onElse) > 0 {
		value += fmt.Sprintf(" %s", onElse)
	}

	query = fmt.Sprintf("(%s %s %s)", constFunctionCase, value, constFunctionEnd)

	if c.alias != nil && len(strings.TrimSpace(*c.alias)) > 0 {
		query += fmt.Sprintf(" %s %s", constFunctionAs, *c.alias)
	}

	return query, nil
}
