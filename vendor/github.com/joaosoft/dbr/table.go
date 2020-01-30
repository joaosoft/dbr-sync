package dbr

type table struct {
	data interface{}
	*functionBase
}

func newTable(db *db, data interface{}) *table {
	return &table{functionBase: newFunctionBase(false, false, db), data: data}
}

func (t *table) Build() (string, error) {
	return handleBuild(t.functionBase, t.data)
}

func (t *table) String() string {
	table, _ := handleExpression(t.functionBase, t.data)
	return table
}
