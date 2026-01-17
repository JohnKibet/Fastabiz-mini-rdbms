package core

type Database interface {
	Exec(query string) (*Result, error)
	Prepare() error
}

// Mirrors real DBs (Exec, Prepare)
