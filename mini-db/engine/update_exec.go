package engine

import (
	"errors"
	"fastabiz-mini-rdbms/mini-db/storage"
)

func (e *Engine) Update(cmd UpdateCommand) (int, error) {
	table, ok := e.Tables[cmd.TableName]
	if !ok {
		return 0, errors.New("table does not exist")
	}

	// Prevent Primary Key update
	if _, exists := cmd.Set[table.PrimaryKey]; exists {
		return 0, errors.New("cannot update primary key")
	}

	if cmd.Where != nil &&
		cmd.Where.Column == table.PrimaryKey &&
		table.PKIndex != nil {
		return e.deleteByPk(table, cmd.Where.Value), nil
	}

	return e.updateByScan(table, cmd), nil
}

func (e *Engine) updateByPK(table *storage.Table, cmd UpdateCommand) int {
  rowID, ok := table.PKIndex.Get(cmd.Where.Value)
  if !ok {
    return 0
  }

  row := table.Rows[storage.RowID(rowID)]
  for col, val := range cmd.Set {
    row[col] = val
  }

  table.Rows[storage.RowID(rowID)] = row
  return 1
}

func (e *Engine) updateByScan(table *storage.Table, cmd UpdateCommand) int {
	updated := 0

	for rowID, row := range table.Rows {
		if row[cmd.Where.Column] == cmd.Where.Value {
			for col, val := range cmd.Set {
				row[col] = val
			}
			table.Rows[rowID] = row
			updated++
		}
	}

	return updated
}