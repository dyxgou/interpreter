package evaluator

import (
	"testing"

	"gihub.com/dyxgou/parser/src/lexer"
	"gihub.com/dyxgou/parser/src/object"
	"gihub.com/dyxgou/parser/src/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	return Eval(program)
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object expected=evaluator.NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}
