package engine

import (
	"fastabiz-mini-rdbms/mini-db/core"
	"fastabiz-mini-rdbms/mini-db/storage"
	"fmt"
	"strings"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (any, error) {
	switch p.current().Type {
	case CREATE:
		return p.parseCreateTable()
	case INSERT:
		return p.parseInsert()
	case SELECT:
		return p.parseSelect()
	case DELETE:
		return p.parseDelete()
	case UPDATE:
		return p.parseUpdate()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.current().Literal)
	}
}

func (p *Parser) parseCreateTable() (*CreateTableCommand, error) {
	p.advance() // CREATE

	_, err := p.expect(TABLE)
	if err != nil {
		return nil, err
	}

	tableNameTok, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(LPAREN)
	if err != nil {
		return nil, err
	}

	var columns []storage.Column

	for {
		// 1. Column name
		colNameTok, err := p.expect(IDENT)
		if err != nil {
			return nil, err
		}

		// 2. Column type
		colTypeTok, err := p.expect(IDENT)
		if err != nil {
			return nil, err
		}

		colType, err := core.ParseDataType(colTypeTok.Literal)
		if err != nil {
			return nil, err
		}

		col := storage.Column{
			Name: colNameTok.Literal,
			Type: colType,
		}

		// 3. PRIMARY KEY
		if p.current().Type == IDENT && strings.ToUpper(p.current().Literal) == "PRIMARY" {
			p.advance() // consume PRIMARY

			keyTok, err := p.expect(IDENT)
			if err != nil {
				return nil, fmt.Errorf("expected KEY after PRIMARY, got %s", keyTok.Literal)
			}

			if strings.ToUpper(keyTok.Literal) != "KEY" {
				return nil, fmt.Errorf("expected KEY after PRIMARY, got %s", keyTok.Literal)
			}

			col.Primary = true
		}

		columns = append(columns, col)

		// 4. Check for comma to continue
		if p.current().Type == COMMA {
			p.advance()
			continue
		}
		break
	}

	_, err = p.expect(RPAREN)
	if err != nil {
		return nil, err
	}

	return &CreateTableCommand{
		TableName: tableNameTok.Literal,
		Columns:   columns,
	}, nil
}


func (p *Parser) parseInsert() (*InsertCommand, error) {
	p.advance() // INSERT
	_, err := p.expect(INTO)
	if err != nil {
		return nil, err
	}

	table, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(LPAREN)
	if err != nil {
		return nil, err
	}

	var cols []string
	for {
		c, err := p.expect(IDENT)
		if err != nil {
			return nil, err
		}
		cols = append(cols, c.Literal)

		if p.current().Type == COMMA {
			p.advance()
			continue
		}
		break
	}

	_, err = p.expect(RPAREN)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(VALUES)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(LPAREN)
	if err != nil {
		return nil, err
	}

	values := make(map[string]any)
	for i := 0; i < len(cols); i++ {
		tok := p.advance()
		values[cols[i]] = tok.Literal

		if p.current().Type == COMMA {
			p.advance()
		}
	}

	_, err = p.expect(RPAREN)
	if err != nil {
		return nil, err
	}

	return &InsertCommand{
		TableName: table.Literal,
		Values:    values,
	}, nil
}

func (p *Parser) parseSelect() (*SelectCommand, error) {
	p.advance() // SELECT

	var cols []string
	for {
		tok := p.advance()
		if tok.Type != IDENT && tok.Type != STAR {
			return nil, fmt.Errorf("invalid column")
		}
		cols = append(cols, tok.Literal)

		if p.current().Type == COMMA {
			p.advance()
			continue
		}
		break
	}

	// FROM
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}

	// main table
	tableTok, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}

	// optional JOIN
	var join *JoinSpec
	if p.current().Type == JOIN {
		join, err = p.parseJoin(tableTok.Literal)
		if err != nil {
			return nil, err
		}
	}

	// optional WHERE
	var where *WhereClause
	if p.current().Type == WHERE {
		p.advance() // WHERE
		col, err := p.expect(IDENT)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(EQ); err != nil {
			return nil, err
		}
		val := p.advance()

		where = &WhereClause{
			Column: col.Literal,
			Value:  val.Literal,
		}
	}

	return &SelectCommand{
		TableName: tableTok.Literal,
		Columns:   cols,
		Join:      join,
		Where:     where,
	}, nil
}

func (p *Parser) parseDelete() (*DeleteCommand, error) {
	p.advance() // DELETE
	_, err := p.expect(FROM)
	if err != nil {
		return nil, err
	}

	table, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}

	p.expect(WHERE)
	col, _ := p.expect(IDENT)
	p.expect(EQ)
	val := p.advance()

	return &DeleteCommand{
		TableName: table.Literal,
		Where: &WhereClause{
			Column: col.Literal,
			Value:  val.Literal,
		},
	}, nil
}

func (p *Parser) parseUpdate() (*UpdateCommand, error) {
	p.advance() // UPDATE
	table, _ := p.expect(IDENT)
	p.expect(SET)

	col, _ := p.expect(IDENT)
	p.expect(EQ)
	val := p.advance()

	p.expect(WHERE)
	wcol, _ := p.expect(IDENT)
	p.expect(EQ)
	wval := p.advance()

	return &UpdateCommand{
		TableName: table.Literal,
		Set: map[string]any{
			col.Literal: val.Literal,
		},
		Where: &WhereClause{
			Column: wcol.Literal,
			Value:  wval.Literal,
		},
	}, nil
}

func (p *Parser) parseJoin(leftTable string) (*JoinSpec, error) {
	p.advance() // JOIN

	// JOIN orders
	rightTableTok, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(ON); err != nil {
		return nil, err
	}

	// users.id
	leftTableTok, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(DOT); err != nil {
		return nil, err
	}
	leftColTok, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(EQ); err != nil {
		return nil, err
	}

	// orders.user_id
	rightTableTok2, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(DOT); err != nil {
		return nil, err
	}
	rightColTok, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}

	// safety check
	if rightTableTok2.Literal != rightTableTok.Literal {
		return nil, fmt.Errorf(
			"JOIN table mismatch: expected %s, got %s",
			rightTableTok.Literal,
			rightTableTok2.Literal,
		)
	}

	return &JoinSpec{
		LeftTable:   leftTableTok.Literal,
		RightTable: rightTableTok.Literal,
		LeftColumn: leftColTok.Literal,
		RightColumn: rightColTok.Literal,
	}, nil
}


func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) expect(t TokenType) (Token, error) {
	tok := p.advance()
	if tok.Type != t {
		return tok, fmt.Errorf("expected %s, got %s", t, tok.Type)
	}
	return tok, nil
}
