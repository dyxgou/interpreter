package parser

import (
	"fmt"
	"testing"

	"gihub.com/dyxgou/parser/src/ast"
)

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp is not an *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value() != value {
		t.Errorf("ident.Value is not %q. got=%q", value, ident.Value())
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral is not %q. got=%q", value, ident.Value())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intli, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if intli.Value != value {
		t.Errorf("intli.Value not %d. got=%d", value, intli.Value)
		return false
	}

	if intli.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("intli.TokenLiteral() not %d. got=%s", value, intli.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.errors

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errors))

	for _, err := range errors {
		t.Error(err)
	}

	t.FailNow()
}

func testBooleanExpression(t *testing.T, exp ast.Expression, expected bool) bool {
	boolExp, ok := exp.(*ast.Boolean)

	if !ok {
		t.Errorf("exp is not *ast.Boolean. got=%T", exp)
		return false
	}

	if boolExp.Value != expected {
		t.Errorf("boolExp.Value is not %t. got=%t", expected, boolExp.Value)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanExpression(t, exp, v)
	}

	t.Errorf("type exp=%T not handeled", exp)
	return false
}

func testInfixExpression(t *testing.T,
	exp ast.Expression,
	left any,
	operator string,
	right any,
) bool {
	infix, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not *ast.InfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, infix.Left, left) {
		t.Log("infix left")
		return false
	}

	if infix.Operator() != operator {
		t.Errorf("infix.Operator expected=%q. got=%q", operator, infix.Operator())
		return false
	}

	if !testLiteralExpression(t, infix.Right, right) {
		t.Log("infix right")
		return false
	}

	return true
}
