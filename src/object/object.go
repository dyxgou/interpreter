package object

import (
	"fmt"
	"strings"

	"github.com/dyxgou/parser/src/ast"
)

type Object interface {
	Type() ObjectType
	String() ObjectString
	Inspect() string
}

type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInFunction
}

func (*BuiltIn) Type() ObjectType { return BuiltInType }
func (*BuiltIn) Inspect() string  { return BuiltInStr }
func (b *BuiltIn) String() string { return b.String() }

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType { return IntegerType }
func (*Integer) Inspect() string  { return IntegerStr }
func (i *Integer) String() string { return fmt.Sprintf("%d", i.Value) }

type String struct {
	Value string
}

func (_ *String) Type() ObjectType { return StringType }
func (_ *String) Inspect() string  { return StringStr }
func (o *String) String() string {
	var sb strings.Builder

	sb.WriteByte('"')
	sb.WriteString(o.Value)
	sb.WriteByte('"')

	return sb.String()
}

type Boolean struct {
	Value bool
}

func (_ *Boolean) Type() ObjectType { return BooleanType }
func (_ *Boolean) Inspect() string  { return BooleanStr }
func (o *Boolean) String() string   { return fmt.Sprintf("%t", o.Value) }

type Null struct{}

func (_ *Null) Type() ObjectType { return NullType }
func (_ *Null) Inspect() string  { return NullStr }
func (_ *Null) String() string   { return NullStr }

type ReturnValue struct {
	Value Object
}

func (_ *ReturnValue) Type() ObjectType { return ReturnType }
func (_ *ReturnValue) Inspect() string  { return ReturnStr }
func (o *ReturnValue) String() string   { return o.Value.Inspect() }

type Error struct {
	Message string
}

func (_ *Error) Type() ObjectType { return ErrorType }
func (_ *Error) Inspect() string  { return ErrorStr }
func (o *Error) String() string   { return fmt.Sprintf("ERROR : %s", o.Message) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Enviroment
}

func (o *Function) Type() ObjectType { return FunctionType }
func (o *Function) Inspect() string  { return FunctionStr }
func (o *Function) String() string {
	var sb strings.Builder

	sb.WriteString("fn")

	sb.WriteByte('(')
	for i, param := range o.Parameters {
		if i > 0 && i < len(o.Parameters) {
			sb.WriteString(" ,")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") { \n\t")

	sb.WriteString(o.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

type Array struct {
	Elements []Object
}

func (*Array) Type() ObjectType      { return ArrayType }
func (*Array) Inspect() ObjectString { return ArrayStr }
func (a *Array) String() string {
	var sb strings.Builder

	sb.WriteByte('[')
	for i, elem := range a.Elements {
		if i > 0 && i < len(a.Elements) {
			sb.WriteString(", ")
		}

		sb.WriteString(elem.String())
	}
	sb.WriteByte(']')

	return sb.String()
}
