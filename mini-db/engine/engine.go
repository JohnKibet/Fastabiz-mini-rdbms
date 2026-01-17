package engine

import (
	"fastabiz-mini-rdbms/mini-db/storage"
)

type Engine struct {
	Tables  map[string]*storage.Table
}

func NewEngine() *Engine {
	return &Engine{
		Tables:  make(map[string]*storage.Table),
	}
}
