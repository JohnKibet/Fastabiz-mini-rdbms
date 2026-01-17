package storage

import "fastabiz-mini-rdbms/mini-db/core"

type Column struct {
	Name    string
	Type    core.DataType
	Primary bool
	Unique  bool
}

// Enforces schema rules
// Enables constraint checks on INSERT
