package evaluator

import (
	"log/slog"
	"testing"

	"github.com/dyxgou/parser/src/lexer"
	"github.com/dyxgou/parser/src/object"
	"github.com/dyxgou/parser/src/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	if p.ErrorsLen() != 0 {
		slog.Error("parser had errors")
	}

	program := p.ParseProgram()
	env := object.NewEnviroment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	intObj, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("the obj is not an *object.Integer. got=%T", obj)
		return false
	}

	if intObj.Value != expected {
		t.Errorf("int.Value expected=%d. got=%d", expected, intObj.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	boolObj, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("obj expected=*object.Boolean. got=%T", obj)
		return false
	}

	if boolObj.Value != expected {
		t.Errorf("boolObj valued expected=%t. got=%t", expected, boolObj.Value)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	strObj, ok := obj.(*object.String)

	if !ok {
		if err, ok := obj.(*object.Error); ok {
			t.Errorf("error evaluating. msg=%q", err.Message)
			return false
		}

		t.Errorf("obj expected=*object.String. got=%T", obj)
		return false
	}

	if strObj.Value != expected {
		t.Errorf("strObj valued expected=%q. got=%q", expected, strObj.Value)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object expected=evaluator.NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}
