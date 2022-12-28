package object

import (
	"fmt"
	"ryanlang/ast"
	"ryanlang/lexer"
	"strconv"
	"strings"
)

type Type int

const (
	_ Type = iota
	ERROR
	NUMBER
	STRING
	FUNCTION
	NULL
	RETURNOBJECT
	CONTINUE
	BREAK
	BOOLEAN
	STRUCT
	MAP
	ARRAY
	MODULE
	EXPORTS
	TUPLE
)

func (t Type) String() string {
	switch t {
	case ERROR:
		return "error"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case FUNCTION:
		return "function"
	case NULL:
		return "null"
	case RETURNOBJECT:
		return "return"
	case BOOLEAN:
		return "boolean"
	case STRUCT:
		return "struct"
	case MAP:
		return "map"
	case ARRAY:
		return "array"
	case MODULE:
		return "module"
	case BREAK:
		return "break"
	case TUPLE:
		return "tuple"
	}

	panic("unknown object type: " + strconv.Itoa(int(t)))
}

type Object interface {
	String() string
	Type() Type
}

func simplehash(obj Object) string {
	return obj.Type().String() + "(" + obj.String() + ")"
}

type Hashable interface {
	Hash() string
}

type Error struct {
	Msg   string
	Child *Error
	Loc   *lexer.Location
}

func (e Error) Type() Type {
	return ERROR
}

func (e Error) String() string {
	result := e.Loc.String() + ": " + e.Msg
	if e.Child != nil {
		result += ": " + e.Child.String()
	}
	return result
}
func (e Error) Last(cnt int) []*Error {
	buf := make([]*Error, cnt)
	cur := &e
	for cur != nil {
		copy(buf[1:], buf)
		buf[0] = cur
		cur = cur.Child
	}
	cut := len(buf)
	for i := range buf {
		if buf[i] == nil {
			cut = i
			break
		}
	}
	return buf[:cut]
}

type Number struct {
	Value int
}

func (n Number) Hash() string {
	return simplehash(n)
}

func (n Number) Type() Type {
	return NUMBER
}

func (n Number) String() string { return strconv.Itoa(n.Value) }

type String struct {
	Value string
}

func (s String) Hash() string {
	return simplehash(s)
}

func (s String) String() string {
	//return "\"" + s.Value + "\""
	return s.Value
}

func (s String) Type() Type {
	return STRING
}

type Function struct {
	Node ast.FuncExpression
	Env  *Environment
}

func (f Function) String() string {
	return f.Node.String()
}

func (f Function) Type() Type {
	return FUNCTION
}

type Null struct{}

func (n Null) String() string {
	return "null"
}

func (n Null) Type() Type {
	return NULL
}

type ReturnObject struct {
	Obj Object
}

func (r ReturnObject) String() string {
	return "return " + r.Obj.String()
}

func (r ReturnObject) Type() Type {
	return RETURNOBJECT
}

type ContinueObject struct{}

func (c ContinueObject) String() string {
	return "continue"
}

func (c ContinueObject) Type() Type {
	return CONTINUE
}

type BreakObject struct{}

func (b BreakObject) String() string {
	return "break"
}
func (b BreakObject) Type() Type {
	return BREAK
}

type Boolean struct {
	Value bool
}

func (b Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b Boolean) Type() Type {
	return BOOLEAN
}

type Struct struct {
	Fields map[string]Object
}

func (s Struct) String() string {
	strs := []string{}
	for k, v := range s.Fields {
		strs = append(strs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("struct{%s}", strings.Join(strs, "; "))
}

func (s Struct) Type() Type {
	return STRUCT
}

type MapItem struct {
	Key   Object
	Value Object
}
type Map struct {
	Fields map[string]MapItem
}

func (m Map) String() string {
	strs := []string{}
	for _, v := range m.Fields {
		strs = append(strs, fmt.Sprintf("%s: %s", v.Key.String(), v.Value.String()))
	}
	return fmt.Sprintf("map{%s}", strings.Join(strs, "; "))
}

func (m Map) Type() Type {
	return MAP
}

type Array struct {
	Items []Object
}

func (a Array) Hash() string {
	strs := []string{}
	for _, item := range a.Items {
		if _, ok := item.(Hashable); ok {
			strs = append(strs, item.(Hashable).Hash())
		} else {
			panic("array contains unhashable type: " + item.Type().String())
		}
	}
	return fmt.Sprintf("%s(%s)", a.Type().String(), strings.Join(strs, ", "))
}

func (a Array) String() string {
	strs := []string{}
	for _, item := range a.Items {
		strs = append(strs, item.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

func (a Array) Type() Type {
	return ARRAY
}

type Module struct {
	Name    string
	Exports map[string]Object
}

func (m Module) String() string {
	var exported string
	if len(m.Exports) > 0 {
		strs := []string{}
		for k := range m.Exports {
			strs = append(strs, k)
		}
		exported = fmt.Sprintf("exports: %s", strings.Join(strs, ", "))
	} else {
		exported = "no exported fields"
	}
	return fmt.Sprintf("(module %s, %s)", m.Name, exported)
}

func (m Module) Type() Type {
	return MODULE
}

type Exports struct {
	Fields []string
}

func (e Exports) String() string {
	return fmt.Sprintf("(exports %s)", strings.Join(e.Fields, ", "))
}

func (e Exports) Type() Type {
	return EXPORTS
}

type Tuple struct {
	Values []Object
}

func (t Tuple) String() string {
	strs := []string{}
	for _, k := range t.Values {
		strs = append(strs, k.String())
	}
	return fmt.Sprintf("tuple(%s)", strings.Join(strs, ", "))
}

func (t Tuple) Type() Type {
	return TUPLE
}
