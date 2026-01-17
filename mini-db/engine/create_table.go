package engine

import (
	"errors"
	"fastabiz-mini-rdbms/mini-db/index"
	"fastabiz-mini-rdbms/mini-db/storage"
)

func (e *Engine) CreateTable(cmd CreateTableCommand) error {
	if _, exists := e.Tables[cmd.TableName]; exists {
		return errors.New("table already exists")
	}

	columnMap := make(map[string]storage.Column)
	var primaryKey string
	var primaryKeyIndex *index.PKIndex

	for _, col := range cmd.Columns {
		if _, exists := columnMap[col.Name]; exists {
			return errors.New("duplicate column: " + col.Name)
		}

		if col.Primary {
			if primaryKey != "" {
				return errors.New("multiple primary keys not allowed")
			}
			primaryKey = col.Name
			primaryKeyIndex = index.NewPKIndex()
		}

		columnMap[col.Name] = col
	}

	if primaryKey == "" {
		return errors.New("primary key required")
	}

	table := &storage.Table{
		Name:       cmd.TableName,
		Columns:    cmd.Columns,
		ColumnMap:  columnMap,
		Rows:       make(map[storage.RowID]storage.Row),
		PrimaryKey: primaryKey,
		AutoInc:    1,
		PKIndex:    primaryKeyIndex,
	}

	e.Tables[cmd.TableName] = table
	return nil

}
