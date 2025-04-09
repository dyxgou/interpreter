package token

type TokenKind byte

const (
	EOF TokenKind = iota
	ILLEGAL

	// Identifies + literals
	IDENT
	STRING
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
	LBRACKET
	RBRACKET

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

var keywords = map[string]TokenKind{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenKind {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
