package object

import (
	"fmt"
)

type ObjectType byte

const (
	IntegerType ObjectType = iota
	BooleanType
	NullType
	ReturnType
	ErrorType
)

type Object interface {
	Type() ObjectType
	String() string
	Inspect() string
}

type Integer struct {
	Value int64
}

func (_ *Integer) Type() ObjectType { return IntegerType }
func (_ *Integer) String() string   { return "INTEGER" }
func (o *Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }

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
