package parser

import (
	"fmt"
	"strconv"

	"gihub.com/dyxgou/parser/src/ast"
	"gihub.com/dyxgou/parser/src/lexer"
	"gihub.com/dyxgou/parser/src/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Precendence byte

const (
	LOWEST Precendence = iota
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precendences = map[token.TokenKind]Precendence{
	token.EQUAL:          EQUALS,
	token.NOT_EQUAL:      EQUALS,
	token.LESS:           LESSGREATER,
	token.GREATER:        LESSGREATER,
	token.PLUS:           SUM,
	token.MINUS:          SUM,
	token.DIVISION:       PRODUCT,
	token.MULTIPLICATION: PRODUCT,
}

type Parser struct {
	l *lexer.Lexer

	errors []error

	curToken  token.Token
	readToken token.Token

	prefixParseFns map[token.TokenKind]prefixParseFn
	infixParseFns  map[token.TokenKind]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFns: make(map[token.TokenKind]prefixParseFn, 20),
		infixParseFns:  make(map[token.TokenKind]infixParseFn, 20),
	}

	// Prefix Funcs
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpresion)
	p.registerPrefix(token.FALSE, p.parseBooleanExpresion)
	p.registerPrefix(token.LPAREN, p.parseGroupingExpression)

	// Infix Funcs
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIVISION, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLICATION, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.GREATER, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.readToken
	p.readToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(k token.TokenKind, fn prefixParseFn) {
	p.prefixParseFns[k] = fn
}

func (p *Parser) registerInfix(k token.TokenKind, fn infixParseFn) {
	p.infixParseFns[k] = fn
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

func (p *Parser) notExpectedTokenErr(expected string, got string) {
	err := fmt.Errorf("expected next token to be '%s' got='%s'", expected, got)

	p.errors = append(p.errors, err)
}

func (p *Parser) notPrefixParseFnError(t token.TokenKind) {
	err := fmt.Errorf("no prefix parse function for %d found", t)
	p.errors = append(p.errors, err)
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) readPrencedence() Precendence {
	if pr, ok := precendences[p.readToken.Kind]; ok {
		return pr
	}

	return LOWEST
}

func (p *Parser) curPrecendence() Precendence {
	if pr, ok := precendences[p.curToken.Kind]; ok {
		return pr
	}

	return LOWEST
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
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.readTokenIs(token.SEMI) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(pr Precendence) ast.Expression {
	prefixFn, ok := p.prefixParseFns[p.curToken.Kind]

	if !ok {
		p.notPrefixParseFnError(p.curToken.Kind)
		return nil
	}

	leftExp := prefixFn()

	for !p.readTokenIs(token.SEMI) && pr < p.readPrencedence() {
		infixFn, ok := p.infixParseFns[p.readToken.Kind]

		if !ok {
			return leftExp
		}

		p.nextToken()

		leftExp = infixFn(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token: p.curToken,
	}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token: p.curToken,
		Left:  left,
	}

	pr := p.curPrecendence()
	p.nextToken()

	exp.Right = p.parseExpression(pr)

	return exp
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectRead(token.IDENT) {
		p.notExpectedTokenErr("variable_name", p.curToken.Literal)
		return nil
	}

	letStmt.Name = &ast.Identifier{Token: p.curToken}

	if !p.expectRead(token.ASSIGN) {
		p.notExpectedTokenErr("=", p.curToken.Literal)
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

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intStmt := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		err := fmt.Errorf("could not parse %s into an Integer", p.curToken.Literal)
		p.errors = append(p.errors, err)
		return nil
	}

	intStmt.Value = v

	return intStmt
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

func (p *Parser) parseBooleanExpresion() ast.Expression {
	// As this function is only going to be called when the curToken.Kind is true or false
	// We can use this function to dont make the conversion manually
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupingExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.readTokenIs(token.RPAREN) {
		return nil
	}

	p.nextToken()
	return exp
}
