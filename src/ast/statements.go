package ast

import "gihub.com/dyxgou/parser/src/token"

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Expression
}

func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }

type ReturnStatement struct {
	Token token.Token
	Expression
}

func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }
