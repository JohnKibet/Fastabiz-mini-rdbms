package engine

import (
	"errors"
	"fastabiz-mini-rdbms/mini-db/storage"
)

func (e *Engine) Join(spec JoinSpec) ([]storage.JoinedRow, error) {
	left, ok := e.Tables[spec.LeftTable]
	if !ok {
		return nil, errors.New("left table not found")
	}

	right, ok := e.Tables[spec.RightTable]
	if !ok {
		return nil, errors.New("right table not found")
	}

	var results []storage.JoinedRow

	for _, lrow := range left.Rows {
		lval := lrow[spec.LeftColumn]

		for _, rrow := range right.Rows {
			rval := rrow[spec.RightColumn]

			if lval == rval {
				merged := make(storage.JoinedRow)

				for col, val := range lrow {
					merged[spec.LeftTable+"."+col] = val
				}
				for col, val := range rrow {
					merged[spec.RightTable+"."+col] = val
				}

				results = append(results, merged)
			}
		}
	}

	return results, nil
}
