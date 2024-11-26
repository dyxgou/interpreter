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

func (e *IntegerLiteral) expressionNode()      {}
func (e *IntegerLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *IntegerLiteral) String() string       { return e.TokenLiteral() }

type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func (e *PrefixExpression) expressionNode()      {}
func (e *PrefixExpression) TokenLiteral() string { return e.Token.Literal }
func (e *PrefixExpression) Operator() string     { return e.TokenLiteral() }

func (e *PrefixExpression) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(e.Operator())
	sb.WriteString(e.Right.String())
	sb.WriteByte(')')

	return sb.String()
}

type InfixExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (e *InfixExpression) expressionNode()      {}
func (e *InfixExpression) TokenLiteral() string { return e.Token.Literal }
func (e *InfixExpression) Operator() string     { return e.TokenLiteral() }

func (e *InfixExpression) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	if e.Left != nil {
		sb.WriteString(e.Left.String())
		sb.WriteByte(' ')
		sb.WriteString(e.Operator())
		sb.WriteByte(' ')
	}

	if e.Right != nil {
		sb.WriteString(e.Right.String())
	}

	sb.WriteByte(')')

	return sb.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (e *Boolean) expressionNode()      {}
func (e *Boolean) TokenLiteral() string { return e.Token.Literal }
func (e *Boolean) String() string       { return e.TokenLiteral() }

type BlockStatement struct {
	Token      token.Token // token.LBRACE "{"
	Statements []Statement
}

func (s *BlockStatement) statementNode()       {}
func (s *BlockStatement) TokenLiteral() string { return s.Token.Literal }
func (s *BlockStatement) String() string {
	var sb strings.Builder

	for _, stmt := range s.Statements {
		sb.WriteString(stmt.String())
	}

	return sb.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (e *IfExpression) expressionNode()      {}
func (e *IfExpression) TokenLiteral() string { return e.Token.Literal }

func (e *IfExpression) String() string {
	var sb strings.Builder

	sb.WriteString("if")
	sb.WriteString(e.Condition.String())
	sb.WriteByte(' ')

	sb.WriteString(e.Consequence.String())
	sb.WriteByte(' ')

	if e.Alternative != nil {
		sb.WriteString("else")
		sb.WriteByte(' ')

		sb.WriteString(e.Alternative.String())
	}

	return sb.String()
}
