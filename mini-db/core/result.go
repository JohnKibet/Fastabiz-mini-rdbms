package core

type Row map[string]any

type Result struct {
	Rows     []Row
	Affected int
	Columns  []string
}
