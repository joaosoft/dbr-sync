package dbr

import "fmt"

type onCase struct {
	value interface{}
}

func newCase(value ...interface{}) *onCase {
	c := &onCase{}
	if len(value) > 0 {
		c.value = value[0]
	}

	return c
}

func (c *onCase) Build(db *db) (_ string, err error) {
	base := newFunctionBase(false, false, db)
	if c.value == nil {
		return "", nil
	}

	var condition string
	condition, err = handleBuild(base, c.value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s ", condition), nil
}
