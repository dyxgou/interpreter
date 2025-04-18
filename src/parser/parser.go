package parser

import (
	"fmt"
	"strconv"

	"github.com/dyxgou/parser/src/ast"
	"github.com/dyxgou/parser/src/lexer"
	"github.com/dyxgou/parser/src/token"
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
	INDEX
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
	token.LPAREN:         CALL,
	token.LBRACKET:       INDEX,
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
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpresion)
	p.registerPrefix(token.FALSE, p.parseBooleanExpresion)
	p.registerPrefix(token.LPAREN, p.parseGroupingExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	// Infix Funcs
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIVISION, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLICATION, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.GREATER, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

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

// returns true and advance the token if the expected token is equal to the current token, if not just returns false and dont advance the token
func (p *Parser) expectRead(k token.TokenKind) bool {
	if p.readTokenIs(k) {
		p.nextToken()
		return true
	}

	return false
}

func (p *Parser) notExpectedTokenErr(expected string, got string) {
	err := fmt.Errorf("expected next token to be %q got=%q", expected, got)

	p.errors = append(p.errors, err)
}

func (p *Parser) notPrefixParseFnError(t token.TokenKind) {
	err := fmt.Errorf("no prefix parse function for %d found", t)
	p.errors = append(p.errors, err)
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) ErrorsLen() int {
	return len(p.errors)
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
		p.notExpectedTokenErr("=", p.readToken.Literal)
		return nil
	}

	p.nextToken()
	letStmt.Value = p.parseExpression(LOWEST)

	if p.readTokenIs(token.SEMI) {
		p.nextToken()
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

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	retStmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	retStmt.Value = p.parseExpression(LOWEST)

	if p.readTokenIs(token.SEMI) {
		p.nextToken()
	}

	return retStmt
}

func (p *Parser) parseBooleanExpresion() ast.Expression {
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

func (p *Parser) parseIfExpression() ast.Expression {
	ifExp := &ast.IfExpression{Token: p.curToken}

	if !p.expectRead(token.LPAREN) {
		p.notExpectedTokenErr("(", p.readToken.Literal)
		return nil
	}

	p.nextToken() // curToken = token.LPAREN
	ifExp.Condition = p.parseExpression(LOWEST)

	if !p.expectRead(token.RPAREN) {
		p.notExpectedTokenErr(")", p.readToken.Literal)
		return nil
	}

	if !p.expectRead(token.LBRACE) {
		p.notExpectedTokenErr("{", p.readToken.Literal)
		return nil
	}

	ifExp.Consequence = p.parseBlockStatement()

	if p.expectRead(token.ELSE) {
		if !p.expectRead(token.LBRACE) {
			p.notExpectedTokenErr("{", p.readToken.Literal)
			return nil
		}

		ifExp.Alternative = p.parseBlockStatement()
	}

	return ifExp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: make([]ast.Statement, 0, 20),
	}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	funcExp := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectRead(token.LPAREN) {
		p.notExpectedTokenErr("(", p.readToken.Literal)
		return nil
	}

	funcExp.Params = p.parseFunctionParams()

	if !p.expectRead(token.LBRACE) {
		p.notExpectedTokenErr("{", p.curToken.Literal)
		return nil
	}

	funcExp.Body = p.parseBlockStatement()

	return funcExp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	al := &ast.ArrayLiteral{Token: p.curToken}

	al.Elements = p.parseExpressionList("]", token.RBRACKET)

	return al
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	ie := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()

	ie.Index = p.parseExpression(LOWEST)
	p.nextToken()

	if !p.curTokenIs(token.RBRACKET) {
		return nil
	}

	return ie
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	params := make([]*ast.Identifier, 0, 20)

	if p.readTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken}
	params = append(params, ident)

	for p.readTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.curToken}
		params = append(params, ident)
	}

	if !p.expectRead(token.RPAREN) {
		p.notExpectedTokenErr(")", p.readToken.Literal)
		return nil
	}

	return params
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	callExp.Arguments = p.parseExpressionList(")", token.RPAREN)

	return callExp
}

func (p *Parser) parseExpressionList(lit string, end token.TokenKind) []ast.Expression {
	elems := make([]ast.Expression, 0, 20)

	if p.expectRead(end) {
		return elems
	}

	p.nextToken()
	elems = append(elems, p.parseExpression(LOWEST))

	for p.readTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		elems = append(elems, p.parseExpression(LOWEST))
	}

	if !p.expectRead(end) {
		p.notExpectedTokenErr(lit, p.readToken.Literal)
		return nil
	}

	return elems
}
