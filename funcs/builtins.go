package funcs

import (
	"bufio"
	"fmt"
	"os"
	"ryanlang/object"
	"strconv"
	"strings"
)

var Builtins = map[string]object.Object{
	"true":  &object.StaticTrue,
	"false": &object.StaticFalse,
	"null":  &object.StaticNull,
}
var BuiltinFunctions = map[string]struct {
	Arguments []string
	Body      func(args map[string]object.Object) object.Object
}{
	"println": {
		Arguments: []string{"s"},
		Body: func(args map[string]object.Object) object.Object {
			fmt.Println(args["s"].String())
			return &object.ReturnObject{Obj: &object.StaticNull}
		},
	},
	"print": {
		Arguments: []string{"s"},
		Body: func(args map[string]object.Object) object.Object {
			fmt.Print(args["s"].String())
			return &object.ReturnObject{Obj: &object.StaticNull}
		},
	},
	"debugger": {
		Arguments: []string{},
		Body: func(args map[string]object.Object) object.Object {
			return &object.ReturnObject{Obj: &object.StaticNull}
		},
	},
	"type": {
		Arguments: []string{"v"},
		Body: func(args map[string]object.Object) object.Object {
			v := args["v"]
			return &object.ReturnObject{Obj: &object.String{Value: v.Type().String()}}
		},
	},
	"len": {
		Arguments: []string{"v"},
		Body: func(args map[string]object.Object) object.Object {
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
		Arguments: []string{"s"},
		Body: func(args map[string]object.Object) object.Object {
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
		Arguments: []string{"i"},
		Body: func(args map[string]object.Object) object.Object {
			i := args["i"]
			if i.Type() != object.NUMBER {
				return &object.Error{Msg: "number parameter expected"}
			}
			s := strconv.Itoa(i.(*object.Number).Value)
			return &object.ReturnObject{Obj: &object.String{Value: s}}
		},
	},
	"strsplit": {
		Arguments: []string{"str", "sep"},
		Body: func(args map[string]object.Object) object.Object {
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
		Arguments: []string{"fn"},
		Body: func(args map[string]object.Object) object.Object {
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
		Arguments: []string{"v"},
		Body: func(args map[string]object.Object) object.Object {
			v := args["v"]
			return &object.ReturnObject{Obj: &object.String{Value: v.String()}}
		},
	},
	"panic": {
		Arguments: []string{"msg"},
		Body: func(args map[string]object.Object) object.Object {
			return &object.Error{Msg: "panic: " + args["msg"].String()}
		},
	},
	"slice": { // todo: replace with a.(from:to) ?
		Arguments: []string{"a", "s", "e"},
		Body: func(args map[string]object.Object) object.Object {
			a := args["a"]
			s := args["s"]
			e := args["e"]
			if a.Type() != object.ARRAY || s.Type() != object.NUMBER || e.Type() != object.NUMBER {
				return &object.Error{Msg: "array and two numbers are expected"}
			}
			return &object.ReturnObject{Obj: &object.Array{Items: append([]object.Object{}, a.(*object.Array).Items[s.(*object.Number).Value:e.(*object.Number).Value]...)}}
		},
	},
	"append": {
		Arguments: []string{"a", "i"},
		Body: func(args map[string]object.Object) object.Object {
			a := args["a"]
			i := args["i"]
			if a.Type() != object.ARRAY {
				return &object.Error{Msg: "array and any value are expected"}
			}

			ret := a.(*object.Array)
			ret.Items = append(ret.Items, i)
			return &object.ReturnObject{Obj: args["a"]}
		},
	},
	"makearray": {
		Arguments: []string{"l", "def"},
		Body: func(args map[string]object.Object) object.Object {
			l := args["l"]
			def := args["def"]
			if l.Type() != object.NUMBER {
				return &object.Error{Msg: "number expected"}
			}
			items := make([]object.Object, l.(*object.Number).Value)
			for i := 0; i < l.(*object.Number).Value; i++ {
				items[i] = def
			}
			return &object.ReturnObject{Obj: &object.Array{Items: items}}
		},
	},
	"has": {
		Arguments: []string{"m", "k"},
		Body: func(args map[string]object.Object) object.Object {
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
		Arguments: []string{"m", "k"},
		Body: func(args map[string]object.Object) object.Object {
			m := args["m"]
			k := args["k"]
			if m.Type() != object.MAP || !object.IsHashable(k) {
				return &object.Error{Msg: "a map and a hashable object expected as arguments"}
			}
			delete(m.(*object.Map).Fields, k.(object.Hashable).Hash())
			return &object.ReturnObject{Obj: &object.StaticNull}
		},
	},
	"iteritems": {
		Arguments: []string{"a"},
		Body: func(args map[string]object.Object) object.Object {
			var ret object.Object = &object.Array{}
			a := args["a"]
			switch a.Type() {
			case object.ARRAY:
				ret.(*object.Array).Items = make([]object.Object, len(a.(*object.Array).Items))
				for i, v := range a.(*object.Array).Items {
					ret.(*object.Array).Items[i] = &object.Array{Items: []object.Object{&object.Number{Value: i}, v}}
				}
			case object.MAP:
				for _, v := range a.(*object.Map).Fields {
					ret.(*object.Array).Items = append(ret.(*object.Array).Items, &object.Array{Items: []object.Object{v.Key, v.Value}})
				}
			default:
				return &object.Error{Msg: "array or map expected, got: " + a.Type().String()}
			}
			return &object.ReturnObject{Obj: ret}
		},
	},
}
