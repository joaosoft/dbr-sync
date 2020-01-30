package dbr

import (
	"fmt"
)

type functionArrayAgg struct {
	value  interface{}
	orders orders

	*functionBase
}

func newFunctionArrayAgg(value interface{}) *functionArrayAgg {
	return &functionArrayAgg{functionBase: newFunctionBase(false, false), value: value}
}

func (c *functionArrayAgg) OrderBy(column string, direction direction) *functionArrayAgg {
	c.orders = append(c.orders, &order{column: column, direction: direction})
	return c
}

func (c *functionArrayAgg) OrderAsc(columns ...string) *functionArrayAgg {
	for _, column := range columns {
		c.orders = append(c.orders, &order{column: column, direction: OrderAsc})
	}
	return c
}

func (c *functionArrayAgg) OrderDesc(columns ...string) *functionArrayAgg {
	for _, column := range columns {
		c.orders = append(c.orders, &order{column: column, direction: OrderDesc})
	}
	return c
}

func (c *functionArrayAgg) Build(db *db) (string, error) {
	c.db = db

	base := newFunctionBase(false, false, db)
	value, err := handleBuild(base, c.value)
	if err != nil {
		return "", err
	}

	orders, err := c.orders.Build()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s(%s%s)", constFunctionArrayAgg, value, orders), nil
}
