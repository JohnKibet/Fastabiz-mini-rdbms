package engine

import (
	"errors"
	"fastabiz-mini-rdbms/mini-db/index"
	"fastabiz-mini-rdbms/mini-db/storage"
)

func (e *Engine) Delete(cmd *DeleteCommand) (int, error) {
	table, ok := e.Tables[cmd.TableName]
	if !ok {
		return 0, errors.New("table does not exist")
	}

	// DELETE without WHERE -> delete all rows
	if cmd.Where == nil {
		count := len(table.Rows)
		table.Rows = make(map[storage.RowID]storage.Row)
		table.NextRowID = 0
		table.PKIndex = index.NewPKIndex()
		return count, nil
	}

	// Fast path: PK-based deletion
	if cmd.Where.Column == table.PrimaryKey && table.PKIndex != nil {
		return e.deleteByPk(table, cmd.Where.Value), nil
	}

	return e.deleteByScan(table, cmd.Where), nil
}

func (e *Engine) deleteByPk(table *storage.Table, value any) int {
	rowID, ok := table.PKIndex.Get(value)
	if !ok {
		return 0
	}

	delete(table.Rows, storage.RowID(rowID))
	table.PKIndex.Delete(value)
	return 1
}

func (e *Engine) deleteByScan(table *storage.Table, where *WhereClause) int {
	deleted := 0

	for rowID, row := range table.Rows {
		if row[where.Column] == where.Value {

			if table.PrimaryKey != "" {
				pkVal := row[table.PrimaryKey]
				table.PKIndex.Delete(pkVal)
			}

			delete(table.Rows, rowID)
			deleted++
		}
	}

	return deleted
}