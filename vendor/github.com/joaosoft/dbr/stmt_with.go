package dbr

import (
	"fmt"
)

type StmtWith struct {
	withs       withs
	isRecursive bool

	dbr         *Dbr
	connections *connections
}

func newStmtWith(dbr *Dbr, connections *connections, name string, isRecursive bool, builder Builder) *StmtWith {
	return &StmtWith{
		dbr:         dbr,
		connections: connections,
		withs:       withs{newWith(name, builder)},
		isRecursive: isRecursive,
	}
}

func (w *StmtWith) With(name string, builder Builder) *StmtWith {
	w.withs = append(w.withs, &with{name: name, builder: builder})

	return w
}

func (w *StmtWith) Select(column ...interface{}) *StmtSelect {
	columns := newColumns(w.connections.Read, false)
	columns.list = column

	return newStmtSelect(w.dbr, w.connections.Read, w, columns)
}

func (w *StmtWith) Insert() *StmtInsert {
	return newStmtInsert(w.dbr, w.connections.Write, w)
}

func (w *StmtWith) Update(table string) *StmtUpdate {
	return newStmtUpdate(w.dbr, w.connections.Write, w, table)
}

func (w *StmtWith) Delete() *StmtDelete {
	return newStmtDelete(w.dbr, w.connections.Write, w)
}

func (w *StmtWith) Build() (string, error) {

	if len(w.withs) == 0 {
		return "", nil
	}

	withs, err := w.withs.Build()
	if err != nil {
		return "", err
	}

	var recursive string
	if w.isRecursive {
		recursive = fmt.Sprintf("%s ", constFunctionRecursive)
	}

	return fmt.Sprintf("%s %s%s", constFunctionWith, recursive, withs), nil
}
