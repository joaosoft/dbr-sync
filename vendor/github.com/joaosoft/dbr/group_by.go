package dbr

import (
	"fmt"
	"strings"
)

type groupBy []string

func (g groupBy) Build() (string, error) {
	if len(g) == 0 {
		return "", nil
	}

	return fmt.Sprintf(" %s %s", constFunctionGroupBy, strings.Join(g, ", ")), nil
}
