package parser

import (
	"testing"

	"github.com/dyxgou/parser/src/ast"
	"github.com/dyxgou/parser/src/lexer"
)

func TestParseStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let five = 5;", "five", 5},
		{"let ten = 10;", "ten", 10},
		{"let theBest = 5614;", "theBest", 5614},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if ps := len(program.Statements); ps != 1 {
			t.Fatalf("program.Statements expected=1 statements. got=%d", ps)
		}

		if !testLetStament(t, program.Statements[0], tt.expectedIdentifier, tt.expectedValue) {
			return
		}
	}
}

func testLetStament(t *testing.T, s ast.Statement, name string, value any) bool {
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

	if !testLiteralExpression(t, letStmt.Value, value) {
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
		value    any
	}{
		{"!5;", "!", 5},
		{"-5;", "-", 5},
		{"!true;", "!", true},
		{"!false;", "!", false},
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

		if !testLiteralExpression(t, prefix.Right, tt.value) {
			return
		}
	}
}

func TestInfixIntExpression(t *testing.T) {
	infixTest := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			t.Log(tt.input)
			return
		}
	}
}

func TestOperatorPrecendence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()

		if actual != tt.expected {
			t.Fatalf("expected=%q. got=%q", tt.expected, actual)
		}
	}
}

func TestParseBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		if ps := len(program.Statements); ps != 1 {
			t.Fatalf("incorrect amount of statements. expected=1. got=%d", ps)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}

		if !testLiteralExpression(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x > y) { x }"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt is not an *ast.ExpressionStatement. got=%T", stmt)
	}

	ifExp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("exp is not an *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, ifExp.Condition, "x", ">", "y") {
		return
	}

	if len(ifExp.Consequence.Statements) != 1 {
		t.Fatalf("consequence expected=1 statements. got=%d", len(ifExp.Consequence.Statements))
	}

	cons, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("cons is not an *ast.ExpressionStatement. got=%T", ifExp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if ifExp.Alternative != nil {
		t.Fatalf("exp.Alternative was not nil. got=%T", ifExp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x > y) { x } else { y }"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if ps := len(program.Statements); ps != 1 {
		t.Fatalf("program.Body expected=1 statements. got=%d", ps)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt is not an *ast.ExpressionStatement. got=%T", stmt)
	}

	ifExp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("exp is not an *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, ifExp.Condition, "x", ">", "y") {
		return
	}

	if len(ifExp.Consequence.Statements) != 1 {
		t.Fatalf("consequence expected=1 statements. got=%d", len(ifExp.Consequence.Statements))
	}

	cons, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("cons is not an *ast.ExpressionStatement. got=%T", ifExp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		t.Log("ifExp.Consequence")
		return
	}

	if ifExp.Alternative == nil {
		t.Fatalf("exp.Alternative was nil.")
	}

	if len(ifExp.Alternative.Statements) != 1 {
		t.Fatalf("alternative expected=1 statements. got=%d", len(ifExp.Alternative.Statements))
	}

	alter, ok := ifExp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alter is not an *ast.ExpressionStatement. got=%T", ifExp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alter.Expression, "y") {
		t.Log("ifExp.Alternative")
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x , y) {x + y;}"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if ps := len(program.Statements); ps != 1 {
		t.Fatalf("program.Statements expected=1 statements. got=%d", ps)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	funcLit, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("exp is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if ps := len(funcLit.Params); ps != 2 {
		t.Fatalf("params expected=2. got=%d", ps)
	}

	testIdentifier(t, funcLit.Params[0], "x")
	testIdentifier(t, funcLit.Params[1], "y")

	if bs := len(funcLit.Body.Statements); bs != 1 {
		t.Fatalf("body expected=1 statements. got=%d", bs)
	}

	exp, ok := funcLit.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("body.Statements not *ast.ExpressionStatement. got=%T", funcLit.Body.Statements[0])
	}

	if !testInfixExpression(t, exp.Expression, "x", "+", "y") {
		return
	}
}

func TestFunctionLiteralParams(t *testing.T) {
	tests := []struct {
		input    string
		params   int
		left     int
		operator string
		right    int
	}{
		{"fn() { 1 + 1; }", 0, 1, "+", 1},
		{"fn(x) { 2 + 1; }", 1, 2, "+", 1},
		{"fn(x, y) { 3 + 8; }", 2, 3, "+", 8},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if ps := len(program.Statements); ps != 1 {
			t.Fatalf("program.Statements expected=1 statements. got=%d", ps)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		funcLit, ok := stmt.Expression.(*ast.FunctionLiteral)

		if !ok {
			t.Fatalf("exp is not *ast.FunctionLiteral. got=%T", stmt.Expression)
		}

		if ps := len(funcLit.Params); ps != tt.params {
			t.Fatalf("params expected=2. got=%d", ps)
		}

		if tt.params != 0 {
			testIdentifier(t, funcLit.Params[0], "x")

			if tt.params == 2 {
				testIdentifier(t, funcLit.Params[1], "y")
			}
		}

		if bs := len(funcLit.Body.Statements); bs != 1 {
			t.Fatalf("body expected=1 statements. got=%d", bs)
		}

		exp, ok := funcLit.Body.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("body.Statements not *ast.ExpressionStatement. got=%T", funcLit.Body.Statements[0])
		}

		if !testInfixExpression(t, exp.Expression, tt.left, tt.operator, tt.right) {
			return
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 + 3, 3 * 3);"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if ps := len(program.Statements); ps != 1 {
		t.Fatalf("program.Statements expected=1 statements. got=%d", ps)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("exp is not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, callExp.Function, "add") {
		return
	}

	if len(callExp.Arguments) != 3 {
		t.Fatalf("call.Arguments expected=2. got=%d", len(callExp.Arguments))
	}

	testLiteralExpression(t, callExp.Arguments[0], 1)
	testInfixExpression(t, callExp.Arguments[1], 2, "+", 3)
	testInfixExpression(t, callExp.Arguments[2], 3, "*", 3)
}

func TestParseStringLiteral(t *testing.T) {
	input := `"hello world";`

	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%t", stmt.Expression)
	}

	if literal.Value() != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value())
	}
}

func TestParseArrayLiteral(t *testing.T) {
	// Empty list works as wel but don't try to test them because for some reason declaring just "[]" gives an error but "let arr = [] is ok"
	input := `[1, 2, 3, 4, 5, 6]`

	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)

	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 6 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
}

func TestParseIndexExpression(t *testing.T) {
	input := "myArray[1 + 1];"

	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	ie, ok := stmt.Expression.(*ast.IndexExpression)

	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, ie.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, ie.Index, 1, "+", 1) {
		return
	}
}
