package object

import "ryanlang/lexer"

func ExpectType(obj Object, typ Type) *Error {
	if obj.Type() != typ {
		return &Error{
			Msg: "expected type: " + typ.String() + ", got: " + obj.Type().String(),
		}
	}

	return nil
}

func WrapError(e Object, msg string, loc *lexer.Location) *Error {
	if !IsError(e) {
		panic("cannot wrap a non-error object")
	}

	return &Error{
		Msg:   msg,
		Child: e.(*Error),
		Loc:   loc,
	}
}

func IsError(obj Object) bool {
	return obj.Type() == ERROR
}
func IsHashable(obj Object) bool {
	_, ok := obj.(Hashable)
	return ok
}
