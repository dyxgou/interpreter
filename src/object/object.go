package object

import "fmt"

type ObjectType byte

const (
	IntegerType ObjectType = iota
	BooleanType
	NullType
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (o *Integer) Type() ObjectType { return IntegerType }
func (o *Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }

type Boolean struct {
	Value bool
}

func (o *Boolean) Type() ObjectType { return BooleanType }
func (o *Boolean) Inspect() string  { return fmt.Sprintf("%t", o.Value) }

type Null struct{}

func (_ *Null) Type() ObjectType { return NullType }
func (_ *Null) Inspect() string  { return "null" }
