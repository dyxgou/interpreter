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
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (_ *Integer) Type() ObjectType { return IntegerType }
func (o *Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }

type Boolean struct {
	Value bool
}

func (_ *Boolean) Type() ObjectType { return BooleanType }
func (o *Boolean) Inspect() string  { return fmt.Sprintf("%t", o.Value) }

type Null struct{}

func (_ *Null) Type() ObjectType { return NullType }
func (_ *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (_ *ReturnValue) Type() ObjectType { return ReturnType }
func (o *ReturnValue) Inspect() string  { return o.Value.Inspect() }
