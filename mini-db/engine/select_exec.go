package engine

import (
	"errors"
	"fastabiz-mini-rdbms/mini-db/storage"
	"fmt"
)

func (e *Engine) Select(cmd SelectCommand) ([]storage.Row, error) {
	table, exists := e.Tables[cmd.TableName]
	if !exists {
		return nil, errors.New("table does not exist")
	}

	var result []storage.Row
	for _, row := range table.Rows {
		// Where filter
		if cmd.Where != nil && cmd.Where.Column == table.PrimaryKey && table.PKIndex != nil {
			rowID, ok := table.PKIndex.Get(cmd.Where.Value)
			if !ok {
				return []storage.Row{}, nil // No matching rows
			}

			row := table.Rows[storage.RowID(rowID)]
			return []storage.Row{row}, nil
		}

		// Projection
		projected := storage.Row{}

		if len(cmd.Columns) == 0 {
			// SELECT *
			for k, v := range row {
				projected[k] = v
			}
		} else {
			for _, col := range cmd.Columns {
				if _, ok := table.ColumnMap[col]; !ok {
					return nil, fmt.Errorf("column %s does not exist in table %s", col, cmd.TableName)
				}
				projected[col] = row[col]
			}
		}

		result = append(result, projected)
	}

	return result, nil
}
