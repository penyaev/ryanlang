package funcs

import (
	"fmt"
	"ryanlang/object"
	"strconv"
)

func expectNoErr(args ...object.Object) object.Object {
	for _, arg := range args {
		if object.IsError(arg) {
			return arg
		}
	}
	return nil
}
func expect(obj object.Object, typ object.Type) object.Object {
	if obj.Type() != typ {
		return object.Error{Msg: fmt.Sprintf("expected: %s, got: %s", typ.String(), obj.Type())}
	}

	return obj
}
func expectHashable(obj object.Object) object.Object {
	if !object.IsHashable(obj) {
		return object.Error{Msg: fmt.Sprintf("expected hashable types, got: %s", obj.Type())}
	}

	return obj
}

func Plus(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value + right.(*object.Number).Value}
	} else if left.Type() == object.STRING && right.Type() == object.STRING {
		return &object.String{Value: left.(*object.String).Value + right.(*object.String).Value}
	} else if left.Type() == object.ARRAY && right.Type() == object.ARRAY {
		items := make([]object.Object, len(left.(*object.Array).Items)+len(right.(*object.Array).Items))
		//for i := 0; i < len(left.(*object.Array).Items); i++ {
		//	items[i] = left.(*object.Array).Items[i]
		//}
		//for i := 0; i < len(right.(*object.Array).Items); i++ {
		//	items[i+len(left.(*object.Array).Items)] = right.(*object.Array).Items[i]
		//}
		copy(items, left.(*object.Array).Items)
		copy(items[len(left.(*object.Array).Items):], right.(*object.Array).Items)
		//ret := append([]object.Object{}, left.(*object.Array).Items...)
		//return &object.Array{Items: append(ret, right.(*object.Array).Items...)}
		return &object.Array{Items: items}

		//return &object.Array{Items: append(left.(*object.Array).Items, right.(*object.Array).Items...)}
	} else { // todo: merge maps?
		return &object.Error{Msg: fmt.Sprintf("incompatible types for the plus operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func Minus(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value - right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for the minus operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func Mult(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value * right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for mult operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func Div(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value / right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for the div operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func Mod(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value % right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for the mod operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func Gt(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return object.StaticBool(left.(*object.Number).Value > right.(*object.Number).Value)
	} else if left.Type() == object.STRING && right.Type() == object.STRING {
		return object.StaticBool(left.(*object.String).Value > right.(*object.String).Value)
	} else {
		return &object.Error{Msg: "don't know how to compare types: " + left.Type().String() + ", " + right.Type().String()} // todo: location
	}
}
func Lt(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return object.StaticBool(left.(*object.Number).Value < right.(*object.Number).Value)
	} else if left.Type() == object.STRING && right.Type() == object.STRING {
		return object.StaticBool(left.(*object.String).Value < right.(*object.String).Value)
	} else {
		return &object.Error{Msg: "don't know how to compare types: " + left.Type().String() + ", " + right.Type().String()} // todo: location
	}
}
func Gte(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return object.StaticBool(left.(*object.Number).Value >= right.(*object.Number).Value)
	} else if left.Type() == object.STRING && right.Type() == object.STRING {
		return object.StaticBool(left.(*object.String).Value >= right.(*object.String).Value)
	} else {
		return &object.Error{Msg: "don't know how to compare types: " + left.Type().String() + ", " + right.Type().String()} // todo: location
	}
}
func Lte(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return object.StaticBool(left.(*object.Number).Value <= right.(*object.Number).Value)
	} else if left.Type() == object.STRING && right.Type() == object.STRING {
		return object.StaticBool(left.(*object.String).Value <= right.(*object.String).Value)
	} else {
		return &object.Error{Msg: "don't know how to compare types: " + left.Type().String() + ", " + right.Type().String()} // todo: location
	}
}
func EqTest(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left.Type() != right.Type() {
		return &object.StaticFalse
	}
	if left.Type() == object.NUMBER {
		return object.StaticBool(left.(*object.Number).Value == right.(*object.Number).Value)
	} else if left.Type() == object.STRING {
		return object.StaticBool(left.(*object.String).Value == right.(*object.String).Value)
	} else if left.Type() == object.BOOLEAN {
		return object.StaticBool(left.(*object.Boolean).Value == right.(*object.Boolean).Value)
	} else if left.Type() == object.NULL {
		return &object.StaticTrue
	}
	return &object.Error{Msg: "don't know how to compare values of type " + left.Type().String()}
}
func FieldAccess(lval object.Object, rval object.Object) object.Object {
	if e := expectNoErr(lval, rval); e != nil {
		return e
	}

	var val object.Object
	if lval.Type() == object.STRUCT {
		if rval = expect(rval, object.STRING); object.IsError(rval) {
			return rval
		}
		var ok bool
		val, ok = lval.(*object.Struct).Fields[rval.(*object.String).Value]
		if !ok {
			return &object.Error{Msg: "field does not exist: " + rval.(*object.String).Value}
		}
	} else if lval.Type() == object.MAP {
		if rval = expectHashable(rval); object.IsError(rval) {
			return rval
		}
		var ok bool
		var item object.MapItem
		item, ok = lval.(*object.Map).Fields[rval.(object.Hashable).Hash()]
		if !ok {
			return &object.Error{Msg: "map item does not exist: " + rval.(object.Hashable).Hash()}
		}
		val = item.Value
	} else if lval.Type() == object.MODULE {
		if rval = expect(rval, object.STRING); object.IsError(rval) {
			return rval
		}
		var ok bool
		val, ok = lval.(*object.Module).Exports[rval.(*object.String).Value]
		if !ok {
			return &object.Error{Msg: "identifier is not exported: " + rval.(*object.String).Value}
		}
	} else if lval.Type() == object.ARRAY {
		if rval = expect(rval, object.NUMBER); object.IsError(rval) {
			return rval
		}
		index := rval.(*object.Number).Value

		if index >= len(lval.(*object.Array).Items) || index < 0 {
			return &object.Error{Msg: "index out of range: " + strconv.Itoa(index)}
		}
		val = lval.(*object.Array).Items[index]
	} else if lval.Type() == object.STRING {
		if rval = expect(rval, object.NUMBER); object.IsError(rval) {
			return rval
		}
		index := rval.(*object.Number).Value

		if index >= len(lval.(*object.String).Value) || index < 0 {
			return &object.Error{Msg: "index out of range: " + strconv.Itoa(index)}
		}
		val = &object.String{Value: string(lval.(*object.String).Value[index])}
	} else {
		return &object.Error{Msg: "field access operator is not supported on this type: " + lval.Type().String()}
	}

	return val
}
func FieldAssign(lval object.Object, rval object.Object, value object.Object) object.Object {
	if e := expectNoErr(lval, rval, value); e != nil {
		return e
	}

	if lval.Type() == object.STRUCT {
		if rval = expect(rval, object.STRING); object.IsError(rval) {
			return rval
		}
		fieldName := rval.(*object.String).Value
		currentValue, ok := lval.(*object.Struct).Fields[fieldName]
		if !ok {
			return &object.Error{
				Msg: "cannot assign to a non-existing field: " + fieldName,
				//Loc: fa.Right.Location(),
			}
		}
		if currentValue.Type() != value.Type() {
			return &object.Error{
				Msg: "field already holds a value of type " + currentValue.Type().String() + ", got: " + value.Type().String(),
			}
		}
		lval.(*object.Struct).Fields[fieldName] = value
	} else if lval.Type() == object.MAP {
		if rval = expectHashable(rval); object.IsError(rval) {
			return rval
		}
		hash := rval.(object.Hashable).Hash()
		currentValue, ok := lval.(*object.Map).Fields[hash]
		if ok && currentValue.Value.Type() != value.Type() {
			return &object.Error{
				Msg: "field already holds a value of type " + currentValue.Value.Type().String() + ", got: " + value.Type().String(),
			}
		}
		lval.(*object.Map).Fields[hash] = object.MapItem{
			Key:   rval,
			Value: value,
		}
	} else if lval.Type() == object.ARRAY {
		if rval = expect(rval, object.NUMBER); object.IsError(rval) {
			return rval
		}
		index := rval.(*object.Number).Value
		if index >= len(lval.(*object.Array).Items) || index < 0 {
			return &object.Error{
				Msg: "cannot assign to an array index out of range: " + strconv.Itoa(index),
			}
		}
		currentValue := lval.(*object.Array).Items[index]
		if currentValue.Type() != value.Type() {
			return &object.Error{
				Msg: "array item already holds a value of type " + currentValue.Type().String() + ", got: " + value.Type().String(),
			}
		}
		lval.(*object.Array).Items[index] = value
	} else {
		return &object.Error{Msg: "unexpected field access assign type: " + lval.Type().String()}
	}

	return value
}
func LogicalOr(left object.Object, right object.Object) object.Object {
	if e := expectNoErr(left, right); e != nil {
		return e
	}

	if left = expect(left, object.BOOLEAN); object.IsError(left) {
		return left
	}
	if right = expect(right, object.BOOLEAN); object.IsError(right) {
		return right
	}

	return object.StaticBool(left.(*object.Boolean).Value || right.(*object.Boolean).Value)
}
