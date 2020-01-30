package dbr

import (
	"fmt"
)

type functionOver struct {
	value     interface{}
	partition *partition
	orders    orders

	*functionBase
}

func newFunctionOver(value interface{}) *functionOver {
	return &functionOver{functionBase: newFunctionBase(false, false), value: value, partition: newPartition()}
}

func (c *functionOver) Partition(column ...interface{}) *functionOver {
	c.partition.list = append(c.partition.list, column...)
	return c
}

func (c *functionOver) OrderBy(column string, direction direction) *functionOver {
	c.orders = append(c.orders, &order{column: column, direction: direction})
	return c
}

func (c *functionOver) OrderAsc(columns ...string) *functionOver {
	for _, column := range columns {
		c.orders = append(c.orders, &order{column: column, direction: OrderAsc})
	}
	return c
}

func (c *functionOver) OrderDesc(columns ...string) *functionOver {
	for _, column := range columns {
		c.orders = append(c.orders, &order{column: column, direction: OrderDesc})
	}
	return c
}

func (c *functionOver) Build(db *db) (string, error) {
	c.db = db

	base := newFunctionBase(false, false, db)
	value, err := handleBuild(base, c.value)
	if err != nil {
		return "", err
	}

	partitions, err := c.partition.Build(db)
	if err != nil {
		return "", err
	}

	orders, err := c.orders.Build()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s(%s%s)", value, constFunctionOver, partitions, orders), nil
}
