package parser

import (
	"fmt"

	"gihub.com/dyxgou/parser/src/ast"
	"gihub.com/dyxgou/parser/src/lexer"
	"gihub.com/dyxgou/parser/src/token"
)

type Parser struct {
	l *lexer.Lexer

	errors []error

	curToken  token.Token
	readToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.readToken
	p.readToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := ast.NewProgram()

	for p.curToken.Kind != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Kind {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectRead(token.IDENT) {
		p.notExpectedTokenErr("variable_name")
		return nil
	}

	letStmt.Name = &ast.Identifier{Token: p.curToken}

	if !p.expectRead(token.ASSIGN) {
		p.notExpectedTokenErr("=")
		return nil
	}

	i := 0
	const scope int = 100

	// TODO: Parse the expression
	for p.curToken.Kind != token.SEMI && i < scope {
		p.nextToken()
		i++
	}

	return letStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	retStmt := &ast.ReturnStatement{Token: p.curToken}

	i := 0
	const scope int = 100

	// TODO: Parse the expression
	for p.curToken.Kind != token.SEMI && i < scope {
		p.nextToken()
		i++
	}

	return retStmt
}

func (p *Parser) readTokenIs(k token.TokenKind) bool {
	if p.readToken.Kind == k {
		return true
	}

	return false
}

func (p *Parser) curTokenIs(k token.TokenKind) bool {
	if p.curToken.Kind == k {
		return true
	}

	return false
}

func (p *Parser) expectRead(k token.TokenKind) bool {
	if p.readTokenIs(k) {
		p.nextToken()
		return true
	}

	return false
}

func (p *Parser) notExpectedTokenErr(expected string) {
	err := fmt.Errorf("expected next token to be '%s' got='%s'", expected, p.readToken.Literal)

	p.errors = append(p.errors, err)
}

func (p *Parser) Errors() []error {
	return p.errors
}
