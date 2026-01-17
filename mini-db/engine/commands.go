package engine

import "fastabiz-mini-rdbms/mini-db/storage"

type CreateTableCommand struct {
	TableName string
	Columns   []storage.Column
}

type InsertCommand struct {
	TableName string
	Values    storage.Row
}

type WhereClause struct {
	Column string
	Value  any
}

type JoinSpec struct {
	LeftTable   string
	RightTable  string
	LeftColumn  string
	RightColumn string
}

type SelectCommand struct {
	TableName string
	Columns   []string
	Join      *JoinSpec
	Where     *WhereClause
}

type DeleteCommand struct {
	TableName string
	Where     *WhereClause
}

type UpdateCommand struct {
	TableName string
	Set       map[string]any
	Where     *WhereClause
}