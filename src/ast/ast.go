package ast

import (
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node

	statementNode()
}

type Expression interface {
	Node

	expressionNode()
}

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{
		Statements: make([]Statement, 0, 20),
	}
}

func (p *Program) String() string {
	var sb strings.Builder

	for _, stmt := range p.Statements {
		sb.WriteString(stmt.String())
	}

	return sb.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}
