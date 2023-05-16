package instruction

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

var ErrTooManyArguments = fmt.Errorf("too many arguments")

type args interface {
	Uint8() (uint8, error)
	Uint16() (uint16, error)
	Uint32() (uint32, error)
}

type argsFromBytes struct {
	b   []byte
	pos int
	cnt int
}

func (a *argsFromBytes) has(l int) bool {
	return len(a.b)-a.pos >= l
}
func (a *argsFromBytes) Uint8() (uint8, error) {
	if !a.has(1) {
		return 0, fmt.Errorf("cannot read uint8 from byte buffer")
	}
	a.cnt++
	a.pos++
	return a.b[a.pos-1], nil
}
func (a *argsFromBytes) Uint16() (uint16, error) {
	if !a.has(2) {
		return 0, fmt.Errorf("cannot read uint16 from byte buffer")
	}
	a.cnt++
	a.pos += 2

	return binary.BigEndian.Uint16(a.b[a.pos-2:]), nil
}
func (a *argsFromBytes) Uint32() (uint32, error) {
	if !a.has(4) {
		return 0, fmt.Errorf("cannot read uint32 from byte buffer")
	}
	a.cnt++
	a.pos += 4

	return binary.BigEndian.Uint32(a.b[a.pos-4:]), nil
}

type argsFromArray struct {
	a   []interface{}
	pos int
}

func (a *argsFromArray) Uint8() (uint8, error) {
	if a.pos >= len(a.a) {
		return 0, fmt.Errorf("cannot fetch uint8 from arg list")
	}
	ret, ok := a.a[a.pos].(int)
	if !ok {
		return 0, fmt.Errorf("expected int, got %s", reflect.TypeOf(a.a[a.pos]).String())
	}
	if ret > math.MaxUint8 {
		return 0, fmt.Errorf("max value of uint8 is %d, got %d", math.MaxUint8, ret)
	}
	a.pos++
	return uint8(ret), nil
}

func (a *argsFromArray) Uint16() (uint16, error) {
	if a.pos >= len(a.a) {
		return 0, fmt.Errorf("cannot fetch uint16 from arg list")
	}
	ret, ok := a.a[a.pos].(int)
	if !ok {
		return 0, fmt.Errorf("expected int, got %s", reflect.TypeOf(a.a[a.pos]).String())
	}
	if ret > math.MaxUint16 {
		return 0, fmt.Errorf("max value of uint16 is %d, got %d", math.MaxUint16, ret)
	}
	a.pos++
	return uint16(ret), nil
}

func (a *argsFromArray) Uint32() (uint32, error) {
	if a.pos >= len(a.a) {
		return 0, fmt.Errorf("cannot fetch uint32 from arg list")
	}
	ret, ok := a.a[a.pos].(int)
	if !ok {
		return 0, fmt.Errorf("expected int, got %s", reflect.TypeOf(a.a[a.pos]).String())
	}
	if ret > math.MaxUint32 {
		return 0, fmt.Errorf("max value of uint32 is %d, got %d", math.MaxUint32, ret)
	}
	a.pos++
	return uint32(ret), nil
}
