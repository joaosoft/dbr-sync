package dbr

import (
	"fmt"
)

type functionExpressions struct {
	comma       bool
	expressions []interface{}

	*functionBase
}

type expression struct {
	value  interface{}
	encode bool
}

func newExpression(value interface{}, encode ...bool) *expression {
	var theEncode bool

	if len(encode) > 0 {
		theEncode = encode[0]
	}
	return &expression{encode: theEncode, value: value}
}

func newFunctionExpressions(comma bool, expressions ...interface{}) *functionExpressions {
	return &functionExpressions{functionBase: newFunctionBase(false, false), expressions: expressions, comma: comma}
}

func (c *functionExpressions) Expression(db *db) (string, error) {
	c.db = db

	if len(c.expressions) == 0 {
		return "", nil
	}

	var expressionValue interface{}
	if value, ok := c.expressions[0].(*expression); ok {
		expressionValue = value.value
	}

	return handleExpression(c.functionBase, expressionValue)
}

func (c *functionExpressions) Build(db *db) (string, error) {
	c.db = db

	var expressions string
	var addComma bool

	lenArgs := len(c.expressions)
	for i, argument := range c.expressions {
		encode := c.functionBase.encode

		if expressionValue, ok := argument.(*expression); ok {
			argument = expressionValue.value
			encode = expressionValue.encode
		}

		expression, err := handleBuild(c.functionBase, argument, encode)
		if err != nil {
			return "", err
		}

		expressions += expression

		if c.comma {
			if expression == constFunctionOpenParentheses {
				addComma = true
				continue
			}

			if expression == constFunctionCloseParentheses {
				addComma = false
				goto next
			}

			if i < lenArgs-1 && c.expressions[i+1] == constFunctionCloseParentheses {
				continue
			}

			if addComma {
				expressions += constFunctionComma
			}
		}

	next:

		if i < lenArgs-1 {
			expressions += " "
		}
	}

	return  fmt.Sprintf("%s", expressions), nil
}
