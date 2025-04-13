package evaluator

import (
	"slices"
	"strings"

	"github.com/dyxgou/parser/src/object"
)

var builtins = map[string]*object.BuiltIn{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if n := len(args); n != 1 {
				return newError("len expected 1 argument. got=%d", n)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to 'len' not supported. got=%T", arg.Inspect())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if n := len(args); n != 1 {
				return newError("function `first` supports just one argument. got=%d", n)
			}

			x := args[0]
			if x.Type() != object.ArrayType {
				return newError("argument to `first` must be an array. got=%q", x.Inspect())
			}

			arr := x.(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if n := len(args); n != 1 {
				return newError("function `last` supports just one argument. got=%d", n)
			}

			x := args[0]
			if x.Type() != object.ArrayType {
				return newError("argument to `last` must be an array. got=%q", x.Inspect())
			}

			arr := x.(*object.Array)
			if n := len(arr.Elements); n > 0 {
				return arr.Elements[n-1]
			}

			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if n := len(args); n != 1 {
				return newError("function `rest` supports just one argument. got=%d", n)
			}

			x := args[0]
			if x.Type() != object.ArrayType {
				return newError("argument to `rest` must be an array. got=%q", x.Inspect())
			}

			arr := x.(*object.Array)
			if n := len(arr.Elements); n > 1 {
				elems := make([]object.Object, n-1, n-1)
				copy(elems, arr.Elements[1:n])

				return &object.Array{Elements: elems}
			}

			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if n := len(args); n != 2 {
				return newError(
					"function `push` supports just two argument. got=%d",
					n,
				)
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError(
					"first argument to `push` must be an array. got=%q (%s)",
					args[0].Inspect(), args[0].String(),
				)
			}

			arr.Elements = append(arr.Elements, args[1])
			return &object.Integer{Value: int64(len(arr.Elements))}
		},
	},
	"pop": {
		Fn: func(args ...object.Object) object.Object {
			if n := len(args); n != 1 {
				return newError("function `pop` supports just one argument. got=%d", n)
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError(
					"first argument to `pop` must be an array. got=%q",
					args[0].Inspect(),
				)
			}

			n := len(arr.Elements)
			item := arr.Elements[n-1]
			arr.Elements = slices.Delete(arr.Elements, n-1, n)

			return item
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var sb strings.Builder

			for _, arg := range args {
				sb.WriteString(arg.String())
			}

			return &object.String{Value: sb.String()}
		},
	},
}

func getBuiltins(name string) (*object.BuiltIn, bool) {
	bi, ok := builtins[name]
	return bi, ok
}
