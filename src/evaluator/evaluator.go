package evaluator

import (
	"gihub.com/dyxgou/parser/src/ast"
	"gihub.com/dyxgou/parser/src/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

const (
	notOperator          string = "!"
	plusOperator         string = "+"
	minusOperator        string = "-"
	productoOperator     string = "*"
	divitionOperator     string = "/"
	greaterOperator      string = ">"
	lessOperator         string = "<"
	greaterEqualOperator string = ">="
	lessEqualOperator    string = "<="
	equalOperator        string = "=="
	notEqualOperator     string = "!="
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		value := Eval(node.Value)
		return &object.ReturnValue{Value: value}
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator(), right)
	case *ast.InfixExpression:
		right := Eval(node.Right)
		left := Eval(node.Left)
		return evalInfixExpression(node.Operator(), right, left)
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		if returnObj, ok := result.(*object.ReturnValue); ok {
			return returnObj.Value
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		if result != nil && result.Type() == object.ReturnType {
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	if right == nil {
		return NULL
	}

	switch operator {
	case notOperator:
		return evalNotOperator(right)
	case minusOperator:
		return evalMinusOperatorExpression(right)
	}

	return NULL
}

func evalNotOperator(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	}

	return FALSE
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerType {
		return NULL
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, right, left object.Object) object.Object {
	if right == nil || left == nil {
		return NULL
	}

	switch {
	case right.Type() == object.IntegerType && left.Type() == object.IntegerType:
		return evalIntegerInfixExpression(operator, right, left)
	case operator == equalOperator:
		return nativeBoolToBooleanObject(left == right)
	case operator == notEqualOperator:
		return nativeBoolToBooleanObject(left != right)
	}

	return NULL
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}

	return FALSE
}

func evalIntegerInfixExpression(operator string, right, left object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case plusOperator:
		return &object.Integer{Value: leftVal + rightVal}
	case minusOperator:
		return &object.Integer{Value: leftVal - rightVal}
	case productoOperator:
		return &object.Integer{Value: leftVal * rightVal}
	case divitionOperator:
		return &object.Integer{Value: leftVal / rightVal}
	case equalOperator:
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case greaterOperator:
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case greaterEqualOperator:
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case lessOperator:
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case lessEqualOperator:
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case notEqualOperator:
		return nativeBoolToBooleanObject(leftVal != rightVal)
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	}

	return true
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}

	return NULL
}
