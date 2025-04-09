package object

import (
	"fmt"
	"strings"

	"github.com/dyxgou/parser/src/ast"
)

type ObjectType byte

const (
	IntegerType ObjectType = iota
	StringType
	BooleanType
	NullType
	ReturnType
	FunctionType
	BuiltInType
	ErrorType
)

type Object interface {
	Type() ObjectType
	String() string
	Inspect() string
}

type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInFunction
}

func (_ *BuiltIn) Type() ObjectType { return BuiltInType }
func (_ *BuiltIn) Inspect() string  { return "built in function" }
func (_ *BuiltIn) String() string   { return "built in function" }

type Integer struct {
	Value int64
}

func (_ *Integer) Type() ObjectType { return IntegerType }
func (_ *Integer) String() string   { return "INTEGER" }
func (o *Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }

type String struct {
	Value string
}

func (_ *String) Type() ObjectType { return StringType }
func (_ *String) String() string   { return "STRING" }
func (o *String) Inspect() string  { return o.Value }

type Boolean struct {
	Value bool
}

func (_ *Boolean) Type() ObjectType { return BooleanType }
func (_ *Boolean) String() string   { return "BOOLEAN" }
func (o *Boolean) Inspect() string  { return fmt.Sprintf("%t", o.Value) }

type Null struct{}

func (_ *Null) Type() ObjectType { return NullType }
func (_ *Null) Inspect() string  { return "null" }
func (_ *Null) String() string   { return "NULL" }

type ReturnValue struct {
	Value Object
}

func (_ *ReturnValue) Type() ObjectType { return ReturnType }
func (_ *ReturnValue) String() string   { return "RETURN" }
func (o *ReturnValue) Inspect() string  { return o.Value.Inspect() }

type Error struct {
	Message string
}

func (_ *Error) Type() ObjectType { return ErrorType }
func (_ *Error) String() string   { return "ERROR" }
func (o *Error) Inspect() string  { return fmt.Sprintf("ERROR : %s", o.Message) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Enviroment
}

func (o *Function) Type() ObjectType { return FunctionType }
func (o *Function) String() string   { return "FUNCTION" }
func (o *Function) Inspect() string {
	var sb strings.Builder

	sb.WriteString("fn")

	sb.WriteByte('(')
	for i, param := range o.Parameters {
		if i > 0 && i < len(o.Parameters) {
			sb.WriteString(" ,")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") { \n")

	sb.WriteString(o.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}
