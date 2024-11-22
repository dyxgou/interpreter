package ast

import (
	"testing"

	"gihub.com/dyxgou/parser/src/token"
)

func TestAstString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Kind:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Kind:    token.IDENT,
						Literal: "myVar",
					},
				},
				Value: &Identifier{
					Token: token.Token{
						Kind:    token.IDENT,
						Literal: "myAnotherVar",
					},
				},
			},
		},
	}

	if ps := program.String(); ps != "let myVar = myAnotherVar;" {
		t.Errorf("program.String() wrong got=%s", ps)
	}
}
