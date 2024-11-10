package lexer

import "gihub.com/dyxgou/parser/src/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	chr          byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

// Modifies the chr field with the current char in the string.
// If the the readPosition is greater than the input, the chr becomes 0.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.chr = 0
	} else {
		l.chr = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() {
	var tok token.Token

	switch l.chr {
	case '=':
		tok = token.New(token.ASSIGN, string(l.chr))
	case ';':
		tok = token.New(token.SEMI, string(l.chr))
	case ',':
		tok = token.New(token.COMMA, string(l.chr))
	}
}
