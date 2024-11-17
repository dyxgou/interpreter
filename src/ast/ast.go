package ast

import "gihub.com/dyxgou/parser/src/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node

	statementNode()
}

type Expression interface {
	Node

	expressionStatement()
}

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{
		Statements: make([]Statement, 0, 20),
	}
}

type Identifier struct {
	Token token.Token // token.IDENT
}

func (i *Identifier) Value() string {
	return i.Token.Literal
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
