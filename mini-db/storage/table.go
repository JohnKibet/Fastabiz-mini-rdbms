package storage

import "fastabiz-mini-rdbms/mini-db/index"

type Table struct {
	Name       string
	Columns    []Column
	ColumnMap  map[string]Column
	Rows       map[RowID]Row
	NextRowID  int


	PrimaryKey string
	PKIndex	 index.Index
	AutoInc    int
}

// map[RowID]Row = fast access
// Indexes mirrors real DB internal catalogs
// AutoInc gives you PK generation cheaply
