package lexer

import (
	"strings"

	"github.com/dyxgou/parser/src/token"
)

type Lexer struct {
	input        string
	position     int // Points to the position of the last read
	readPosition int // Points to the reading position
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
	}

	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = byte(token.EOF)
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	switch l.ch {
	default:
		if isLetter(l.ch) {
			t.Literal = l.readIdentifier()
			t.Kind = token.LookupIdent(t.Literal)
			return t
		} else if isDigit(l.ch) {
			t.Literal = l.readNumber()
			t.Kind = token.INT
			return t
		} else {
			t = token.New(token.ILLEGAL, string(l.ch))
		}
	case '=':
		if ch := l.peekChar(); ch == '=' {
			t = token.New(token.EQUAL, getCompositeString(l.ch, ch))
			l.readChar()
			break
		}
		t = token.New(token.ASSIGN, string(l.ch))
	case '{':
		t = token.New(token.LBRACE, string(l.ch))
	case '}':
		t = token.New(token.RBRACE, string(l.ch))
	case '(':
		t = token.New(token.LPAREN, string(l.ch))
	case ')':
		t = token.New(token.RPAREN, string(l.ch))
	case ',':
		t = token.New(token.COMMA, string(l.ch))
	case ';':
		t = token.New(token.SEMI, string(l.ch))
	case '+':
		t = token.New(token.PLUS, string(l.ch))
	case '-':
		t = token.New(token.MINUS, string(l.ch))
	case '*':
		t = token.New(token.MULTIPLICATION, string(l.ch))
	case '/':
		t = token.New(token.DIVISION, string(l.ch))
	case '<':
		t = token.New(token.LESS, string(l.ch))
	case '>':
		t = token.New(token.GREATER, string(l.ch))
	case '"':
		t = token.New(token.QOUTE, string(l.ch))
	case '!':
		if ch := l.peekChar(); ch == '=' {
			t = token.New(token.NOT_EQUAL, getCompositeString(l.ch, ch))
			l.readChar()
			break
		}
		t = token.New(token.NOT, string(l.ch))
	case byte(token.EOF):
		t = token.New(token.EOF, "")
	}

	l.readChar()
	return t
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	pos := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return byte(token.EOF)
	}

	return l.input[l.readPosition]
}

func getCompositeString(b ...byte) string {
	var s strings.Builder

	s.Write(b)

	return s.String()
}
