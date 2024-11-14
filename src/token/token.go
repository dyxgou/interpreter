package token

type TokenKind byte

const (
	EOF TokenKind = iota
	ILLEGAL

	// Identifies + literals
	IDENT
	INT

	// Operators
	ASSIGN
	PLUS
	MINUS
	MULTIPLICATION
	DIVISION
	LESS
	GREATER
	NOT
	EQUAL
	NOT_EQUAL

	// Delimiters
	COMMA
	SEMI

	LPAREN
	RPAREN
	LBRACE
	RBRACE
	QOUTE

	// Keywords
	FUNCTION
	LET
	TRUE
	FALSE
	RETURN
	IF
	ELSE
)

type Token struct {
	Kind    TokenKind
	Literal string
}

// Creates a new token
func New(k TokenKind, literal string) Token {
	return Token{
		Kind:    k,
		Literal: literal,
	}
}

