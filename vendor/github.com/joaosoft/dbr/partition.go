package dbr

import "fmt"

type partition struct {
	*columns
}

func newPartition() *partition {
	return &partition{columns: newColumns(nil, false)}
}

func (p *partition) Build(db *db) (string, error) {
	p.db = db

	value, err := p.columns.Build()
	if err != nil {
		return "", err
	}

	if value == "" {
		return "", nil
	}

	return fmt.Sprintf("%s %s", constFunctionPartitionBy, value), nil
}
