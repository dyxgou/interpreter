package lexer

import (
	"testing"

	"gihub.com/dyxgou/parser/src/token"
)

func TestReadPosition(t *testing.T) {
	input := "=(){}"

	l := New(input)

	for l.ch != byte(token.EOF) {
		l.readChar()
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
x + y;
};
let result = add(five, ten);
!-/*+5;
5 < 10 > 5;

if (7 < 10) {
  return true;
} else {
  return false;
}
`
	tests := []struct {
		expectedKind    token.TokenKind
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMI, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMI, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMI, ";"},
		{token.RBRACE, "}"},
		{token.SEMI, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMI, ";"},
		{token.NOT, "!"},
		{token.MINUS, "-"},
		{token.DIVISION, "/"},
		{token.MULTIPLICATION, "*"},
		{token.PLUS, "+"},
		{token.INT, "5"},
		{token.SEMI, ";"},
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.GREATER, ">"},
		{token.INT, "5"},
		{token.SEMI, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "7"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMI, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMI, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Kind != tt.expectedKind {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedKind, tok.Kind)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestDoubleToken(t *testing.T) {
	input := `10 == 10; 
  9 != 10;`

	tests := []struct {
		expectedKind    token.TokenKind
		expectedLiteral string
	}{
		{expectedKind: token.INT, expectedLiteral: "10"},
		{expectedKind: token.EQUAL, expectedLiteral: "=="},
		{expectedKind: token.INT, expectedLiteral: "10"},
		{expectedKind: token.SEMI, expectedLiteral: ";"},
		{expectedKind: token.INT, expectedLiteral: "9"},
		{expectedKind: token.NOT_EQUAL, expectedLiteral: "!="},
		{expectedKind: token.INT, expectedLiteral: "10"},
		{expectedKind: token.SEMI, expectedLiteral: ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		t.Log(tok)

		if tok.Kind != tt.expectedKind {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d",

				i, tt.expectedKind, tok.Kind)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
