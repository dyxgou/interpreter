package parser

import (
	"fmt"
	"testing"

	"gihub.com/dyxgou/parser/src/ast"
	"gihub.com/dyxgou/parser/src/lexer"
)

func TestParseStatements(t *testing.T) {
	input := `let five = 5;
  let ten = 10;
  let foobar = 1221;`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatal("Program returned nil")
	}

	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatal("The program doesnt have the exact staments")
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"five"},
		{"ten"},
		{"foobar"},
	}

	for i, stmt := range program.Statements {
		tt := tests[i]

		if !testLetStament(t, stmt, tt.expectedIdentifier) {
			t.Fail()
		}
	}
}

func testLetStament(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement got=%T", s)
		return false
	}

	if letStmt.Name.Value() != name {
		t.Errorf("letStmt.Name.Value not '%s' got='%s'", name, letStmt.Name.Value())
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s' got='%s'", name, letStmt.Name.Value())
		return false
	}

	return true
}

func TestFailingProgramParsing(t *testing.T) {
	input := "let = 5;"

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatal("Program returned nil")
	}

	if len(p.errors) <= 0 {
		t.Fatal("the parser hasn't detected all the errors")
	}

}

func TestReturnStatement(t *testing.T) {
	input := `return add(123 + 123 - 321);
  return 5;
  return 4;`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatal("Program returned nil")
	}

	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatal("The program doesnt have the exact staments")
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Fatalf("return statement expected 'ReturnStatement' got='%T'", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("return statement kind expected 'return' got='%s'", returnStmt.TokenLiteral())
		}
	}

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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if stmts := len(program.Statements); stmts != 1 {
		t.Fatalf("the program has not enoght statements got=%d", stmts)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", stmt)
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%s", ident.Value())
	}

	if tl := ident.TokenLiteral(); tl != "foobar" {
		t.Fatalf("ident.TokenLiteral() not=%s. got=%s", "foobar", tl)
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if ps := len(program.Statements); ps != 1 {
		t.Fatalf("program has not enoght statements. got=%d", ps)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.NumberExpression. got=%T", stmt)
	}

	if !testIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTest := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-5;", "-", 5},
	}

	for _, tt := range prefixTest {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if ps := len(program.Statements); ps != 1 {
			t.Log(program.Statements)
			t.Fatalf("program has not enoght statements. got=%d", ps)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not an *ast.NumberExpression. got=%T", stmt)
		}

		prefix, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt.Expresion is not an *ast.PrefixExpression. got=%T", stmt)
		}

		if ope := prefix.Operator(); ope != tt.operator {
			t.Fatalf("prefix.Operator is not '%s'. got=%s", tt.operator, ope)
		}

		if !testIntegerLiteral(t, prefix.Right, tt.value) {
			return
		}
	}
}

func TestInfixIntExpression(t *testing.T) {
	infixTest := []struct {
		input      string
		rightValue int64
		operator   string
		leftValue  int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTest {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if ps := len(program.Statements); ps != 1 {
			t.Fatalf("program expected=1 statements. got=%d", ps)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not an *ast.NumberExpression. got=%T", stmt)
		}

		infix, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, infix.Left, tt.leftValue) {
			t.Log("infix.Left")
			return
		}

		if infix.Operator() != tt.operator {
			t.Fatalf("exp.Operator expected=%s. got=%s", tt.operator, infix.Operator())
		}

		if !testIntegerLiteral(t, infix.Right, tt.rightValue) {
			t.Log("infix.Right")
			return
		}
	}
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
