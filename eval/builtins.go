package eval

import (
	"bufio"
	"fmt"
	"os"
	"ryanlang/ast"
	"ryanlang/object"
	"strconv"
	"strings"
)

var builtins = map[string]object.Object{
	"true":  &object.Boolean{Value: true},
	"false": &object.Boolean{Value: false},
}
var builtinFunctions = map[string]struct {
	arguments []string
	body      func(env *object.Environment, args map[string]object.Object) object.Object
}{
	"println": {
		arguments: []string{"s"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			fmt.Println(args["s"].String())
			return &object.ReturnObject{Obj: &object.Null{}}
		},
	},
	"print": {
		arguments: []string{"s"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			fmt.Print(args["s"].String())
			return &object.ReturnObject{Obj: &object.Null{}}
		},
	},
	"type": {
		arguments: []string{"v"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			v := args["v"]
			return &object.ReturnObject{Obj: &object.String{Value: v.Type().String()}}
		},
	},
	"len": {
		arguments: []string{"v"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			v := args["v"]

			var val int
			switch v.Type() {
			case object.STRING:
				val = len(v.(*object.String).Value)
			case object.ARRAY:
				val = len(v.(*object.Array).Items)
			case object.MAP:
				val = len(v.(*object.Map).Fields)
			default:
				return &object.Error{Msg: "cannot calculate len() on type " + v.Type().String()}
			}

			return &object.ReturnObject{Obj: &object.Number{Value: val}}
		},
	},
	"atoi": {
		arguments: []string{"s"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			s := args["s"]
			if s.Type() != object.STRING {
				return &object.Error{Msg: "string parameter expected"}
			}
			i, err := strconv.Atoi(s.(*object.String).Value)
			if err != nil {
				return &object.Error{Msg: err.Error()}
			}
			return &object.ReturnObject{Obj: &object.Number{Value: i}}
		},
	},
	"itoa": {
		arguments: []string{"i"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			i := args["i"]
			if i.Type() != object.NUMBER {
				return &object.Error{Msg: "number parameter expected"}
			}
			s := strconv.Itoa(i.(*object.Number).Value)
			return &object.ReturnObject{Obj: &object.String{Value: s}}
		},
	},
	"strsplit": {
		arguments: []string{"str", "sep"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			str := args["str"]
			sep := args["sep"]
			if str.Type() != object.STRING || sep.Type() != object.STRING {
				return &object.Error{Msg: "string parameters expected"}
			}
			res := strings.Split(str.(*object.String).Value, sep.(*object.String).Value)
			ret := &object.Array{}
			for _, s := range res {
				ret.Items = append(ret.Items, &object.String{Value: s})
			}
			return &object.ReturnObject{Obj: ret}
		},
	},
	"readlines": {
		arguments: []string{"fn"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			fn := args["fn"]
			if fn.Type() != object.STRING {
				return &object.Error{Msg: "string parameter expected"}
			}

			f, err := os.Open(fn.(*object.String).Value)
			if err != nil {
				return &object.Error{Msg: err.Error()}
			}
			ret := []object.Object{}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				ret = append(ret, &object.String{Value: scanner.Text()})
			}
			if err := scanner.Err(); err != nil {
				return &object.Error{Msg: err.Error()}
			}
			return &object.ReturnObject{Obj: &object.Array{Items: ret}}
		},
	},
	"dump": {
		arguments: []string{"v"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			v := args["v"]
			return &object.ReturnObject{Obj: &object.String{Value: v.String()}}
		},
	},
	"panic": {
		arguments: []string{"msg"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			return &object.Error{Msg: "panic: " + args["msg"].String()}
		},
	},
	"slice": { // todo: replace with a.(from:to) ?
		arguments: []string{"a", "s", "e"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			a := args["a"]
			s := args["s"]
			e := args["e"]
			if a.Type() != object.ARRAY || s.Type() != object.NUMBER || e.Type() != object.NUMBER {
				return &object.Error{Msg: "array and two numbers are expected"}
			}
			return &object.ReturnObject{Obj: &object.Array{Items: append([]object.Object{}, a.(*object.Array).Items[s.(*object.Number).Value:e.(*object.Number).Value]...)}}
		},
	},
	"range": {
		arguments: []string{"l"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			l := args["l"]
			if l.Type() != object.NUMBER {
				return &object.Error{Msg: "number expected"}
			}
			items := make([]object.Object, l.(*object.Number).Value)
			for i := 0; i < l.(*object.Number).Value; i++ {
				items[i] = &object.Number{Value: i}
			}
			return &object.ReturnObject{Obj: &object.Array{Items: items}}
		},
	},
	"has": {
		arguments: []string{"m", "k"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			m := args["m"]
			k := args["k"]
			if m.Type() != object.MAP || !object.IsHashable(k) {
				return &object.Error{Msg: "a map and a hashable object expected as arguments"}
			}
			_, ok := m.(*object.Map).Fields[k.(object.Hashable).Hash()]
			return &object.ReturnObject{Obj: &object.Boolean{Value: ok}}
		},
	},
	"delete": {
		arguments: []string{"m", "k"},
		body: func(env *object.Environment, args map[string]object.Object) object.Object {
			m := args["m"]
			k := args["k"]
			if m.Type() != object.MAP || !object.IsHashable(k) {
				return &object.Error{Msg: "a map and a hashable object expected as arguments"}
			}
			delete(m.(*object.Map).Fields, k.(object.Hashable).Hash())
			return &object.ReturnObject{Obj: &object.Null{}}
		},
	},
}

func newBuiltinEnvironment() *object.Environment {
	env := object.NewEnvironment()

	for key, value := range builtinFunctions {
		var args []ast.Identifier
		for _, arg := range value.arguments {
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

	for key, value := range builtins {
		env.Set(key, value)
	}

	return env
}
