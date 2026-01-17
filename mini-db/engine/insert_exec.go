package engine

import (
	"errors"
	"fastabiz-mini-rdbms/mini-db/storage"
)

func (e *Engine) Insert(cmd InsertCommand) error {
	table, ok := e.Tables[cmd.TableName]
	if !ok {
		return errors.New("table does not exist")
	}

	row := storage.Row{}
	for col, val := range cmd.Values {
		row[col] = val
	}

	rowID := table.NextRowID

	// Primary Key enforcement
	if table.PrimaryKey != "" {
		pkVal, ok := row[table.PrimaryKey]
		if !ok {
			return errors.New("primary key missing")
		}
		if err := table.PKIndex.Insert(pkVal, rowID); err != nil {
			return err
		}
	}

	table.Rows[storage.RowID(rowID)] = row
	table.NextRowID++
	
	return nil
}
