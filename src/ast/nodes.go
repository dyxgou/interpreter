package ast

import (
	"strings"

	"gihub.com/dyxgou/parser/src/token"
)

type Identifier struct {
	Token token.Token // token.IDENT
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Value() string {
	return i.Token.Literal
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.TokenLiteral()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }
func (s *LetStatement) String() string {
	var sb strings.Builder

	sb.WriteString(s.TokenLiteral() + " ")
	sb.WriteString(s.Name.String())
	sb.WriteString(" = ")

	if s.Value != nil {
		sb.WriteString(s.Value.String())
	}

	sb.WriteString(";")

	return sb.String()
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ReturnStatement) String() string {
	var sb strings.Builder

	sb.WriteString(s.TokenLiteral() + " ")

	if s.Value != nil {
		sb.WriteString(s.Value.String())
	}

	sb.WriteString(";")

	return sb.String()
}

type ExpressionStatement struct {
	Token token.Token
	Expression
}

func (s *ExpressionStatement) statementNode()       {}
func (s *ExpressionStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ExpressionStatement) String() string {
	var sb strings.Builder

	if s.Expression != nil {
		sb.WriteString(s.Expression.String())
	}

	return sb.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (s *IntegerLiteral) expressionNode()      {}
func (s *IntegerLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *IntegerLiteral) String() string       { return s.TokenLiteral() }

type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func (s *PrefixExpression) expressionNode()  {}
func (s *PrefixExpression) Operator() string { return s.Token.Literal }

func (s *PrefixExpression) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(s.Operator())
	sb.WriteString(s.Right.String())
	sb.WriteByte(')')

	return sb.String()
}

func (s *PrefixExpression) TokenLiteral() string {
	return s.String()
}
