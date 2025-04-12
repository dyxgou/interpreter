package object

type ObjectType byte
type ObjectString = string

const (
	IntegerType ObjectType = iota
	StringType
	BooleanType
	NullType
	ReturnType
	FunctionType
	BuiltInType
	ArrayType
	ErrorType
)

const (
	IntegerStr  ObjectString = "INTEGER"
	StringStr   ObjectString = "STRING"
	BooleanStr  ObjectString = "BOOLEAN"
	NullStr     ObjectString = "NULL"
	ReturnStr   ObjectString = "RETURN"
	FunctionStr ObjectString = "FUNCTION"
	BuiltInStr  ObjectString = "BUILTIN"
	ArrayStr    ObjectString = "ARRAY"
	ErrorStr    ObjectString = "ERROR"
)
