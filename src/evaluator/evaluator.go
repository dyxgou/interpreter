package evaluator

import (
	"fmt"

	"github.com/dyxgou/parser/src/ast"
	"github.com/dyxgou/parser/src/object"
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

func Eval(node ast.Node, env *object.Enviroment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}
		env.Set(node.Name.Value(), val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.ReturnStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator(), right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
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

func evalProgram(stmts []ast.Statement, env *object.Enviroment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement, env *object.Enviroment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if result != nil {
			if result.Type() == object.ReturnType || result.Type() == object.ErrorType {
				return result
			}
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

	return newError(
		fmt.Sprintf("unknown operator: %s%s", operator, right.String()),
	)
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
		return newError(
			fmt.Sprintf("unknown operator: %s%s", "-", right.String()),
		)
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, right, left object.Object) object.Object {
	if right == nil || left == nil {
		return NULL
	}

	switch {
	case right.Type() != left.Type():
		return newError("type mismatch: %s %s %s", left.String(), operator, right.String())
	case right.Type() == object.IntegerType && left.Type() == object.IntegerType:
		return evalIntegerInfixExpression(operator, right, left)
	case operator == equalOperator:
		return nativeBoolToBooleanObject(left == right)
	case operator == notEqualOperator:
		return nativeBoolToBooleanObject(left != right)
	}

	return newError("unknown operator: %s %s %s", right.String(), operator, left.String())
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

func evalIfExpression(ie *ast.IfExpression, env *object.Enviroment) object.Object {
	condition := Eval(ie.Condition, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func evalIdentifier(node *ast.Identifier, env *object.Enviroment) object.Object {
	val, ok := env.Get(node.Value())

	if !ok {
		return newError("identifier not found: %s", node.Value())
	}

	return val
}

func newError(message string, a ...any) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(message, a...),
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorType
	}

	return false
}
