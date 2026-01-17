package engine

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

type Tokenizer struct {
	input string
	pos   int
	ch    byte
}

const (
	// Special
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	// Operators
	EQ     TokenType = "="
	COMMA  TokenType = ","
	STAR   TokenType = "*"
	LPAREN TokenType = "("
	RPAREN TokenType = ")"

	// Keywords
	SELECT TokenType = "SELECT"
	FROM   TokenType = "FROM"
	WHERE  TokenType = "WHERE"
	INSERT TokenType = "INSERT"
	INTO   TokenType = "INTO"
	VALUES TokenType = "VALUES"
	CREATE TokenType = "CREATE"
	TABLE  TokenType = "TABLE"
	DELETE TokenType = "DELETE"
	UPDATE TokenType = "UPDATE"
	SET    TokenType = "SET"
	JOIN   TokenType = "JOIN"
	ON     TokenType = "ON"
	DOT    TokenType = "DOT"
)

var keywords = map[string]TokenType{
	"select": SELECT,
	"from":   FROM,
	"where":  WHERE,
	"insert": INSERT,
	"into":   INTO,
	"values": VALUES,
	"create": CREATE,
	"table":  TABLE,
	"delete": DELETE,
	"update": UPDATE,
	"set":    SET,
	"join":   JOIN,
	"on":     ON,
	"dot":    DOT,
}

func NewTokenizer(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.readChar()
	return t
}

func Tokenize(input string) ([]Token, error) {
	t := NewTokenizer(input)
	var tokens []Token

	for {
		tok := t.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == ILLEGAL {
			return nil, fmt.Errorf("illegal token: %s", tok.Literal)
		}

		if tok.Type == EOF {
			break
		}
	}

	return tokens, nil
}


func (t *Tokenizer) NextToken() Token {
	t.skipWhitespace()

	switch t.ch {
	case '=':
		tok := Token{Type: EQ, Literal: "="}
		t.readChar()
		return tok
	case ',':
		tok := Token{Type: COMMA, Literal: ","}
		t.readChar()
		return tok
	case '*':
		tok := Token{Type: STAR, Literal: "*"}
		t.readChar()
		return tok
	case '(':
		tok := Token{Type: LPAREN, Literal: "("}
		t.readChar()
		return tok
	case ')':
		tok := Token{Type: RPAREN, Literal: ")"}
		t.readChar()
		return tok
	case '\'':
		return t.readString()
	case 0:
		return Token{Type: EOF}
	default:
		if isLetter(t.ch) {
			literal := t.readIdentifier()
			tokenType := lookupIdent(literal)
			return Token{Type: tokenType, Literal: literal}
		} else if isDigit(t.ch) {
			return Token{Type: NUMBER, Literal: t.readNumber()}
		}

		tok := Token{Type: ILLEGAL, Literal: string(t.ch)}
		t.readChar()
		return tok
	}
}

func (t *Tokenizer) readChar() {
	if t.pos >= len(t.input) {
		t.ch = 0
	} else {
		t.ch = t.input[t.pos]
	}
	t.pos++
}

func (t *Tokenizer) skipWhitespace() {
	for t.ch == ' ' || t.ch == '\n' || t.ch == '\t' || t.ch == '\r' {
		t.readChar()
	}
}

func (t *Tokenizer) readIdentifier() string {
	start := t.pos - 1
	for isLetter(t.ch) || isDigit(t.ch) {
		t.readChar()
	}
	return t.input[start : t.pos-1]
}

func (t *Tokenizer) readNumber() string {
	start := t.pos - 1
	for isDigit(t.ch) {
		t.readChar()
	}
	return t.input[start : t.pos-1]
}

func (t *Tokenizer) readString() Token {
	t.readChar() // skip opening quote
	start := t.pos - 1

	for t.ch != '\'' && t.ch != 0 {
		t.readChar()
	}

	lit := t.input[start : t.pos-1]
	t.readChar() // skip closing quote

	return Token{Type: STRING, Literal: lit}
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
