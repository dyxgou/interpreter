package token

type TokenKind byte

const (
	ILLEGAL TokenKind = iota
	EOF

	// Identifies + literals
	IDENT
	INT

	// Operators
	ASSIGN
	PLUS

	// Delimiters
	COMMA
	SEMI

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	// Keywords
	FUNCTION
	LET
)

type Token struct {
	kind    TokenKind
	Literal string
}

// Creates a new token
func New(k TokenKind, literal string) Token {
	return Token{
		kind:    k,
		Literal: literal,
	}
}
