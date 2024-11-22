package parser

import (
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

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if literal.Value != 5 {
		t.Fatalf("stmt.Value is not 5. got=%d", literal.Value)
	}

	if tl := literal.TokenLiteral(); tl != "5" {
		t.Fatalf("stmt.TokenLiteral() is not '5'. got=%s", tl)
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTest := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"!5;", "!", 5},
	}

	for _, tt := range prefixTest {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		if ps := len(program.Statements); ps != 1 {
			t.Fatalf("program has not enoght statements. got=%d", ps)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not an *ast.NumberExpression. got=%T", stmt)
		}

		prefix, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {

		}

	}

}
