package evaluator

import "github.com/dyxgou/parser/src/object"

var builtins = map[string]*object.BuiltIn{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if n := len(args); n != 1 {
				return newError("len expected 1 argument. got=%d", n)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to 'len' not supoorted. got=%T", arg.Inspect())
			}
		},
	},
}
