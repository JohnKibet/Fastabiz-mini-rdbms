package main

import (
	"fastabiz-mini-rdbms/mini-db/engine"
	"fastabiz-mini-rdbms/mini-db/repl"
)

func main() {
	db := engine.NewEngine()
	repl.New(db).Run()
}
