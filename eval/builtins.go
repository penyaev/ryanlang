package eval

import (
	"ryanlang/ast"
	"ryanlang/funcs"
	"ryanlang/object"
)

func newBuiltinEnvironment() *object.Environment {
	env := object.NewEnvironment()

	for key, value := range funcs.BuiltinFunctions {
		var args []ast.Identifier
		for _, arg := range value.Arguments {
			args = append(args, ast.Identifier{Name: arg})
		}
		env.Set(key, &object.Function{
			Env: env,
			Node: ast.FuncExpression{
				Arguments: args,
				Body:      ast.BuiltinFunction{Name: key},
			},
		})
	}

	for key, value := range funcs.Builtins {
		env.Set(key, value)
	}

	return env
}
