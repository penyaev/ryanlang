package instruction

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

func Size(op Op) int {
	switch op {
	case OpAdd, OpMult, OpGt, OpGte, OpLt, OpLte, OpClosure, OpSub, OpDiv, OpMod, OpEqTest, OpFieldAccess, OpFieldAssign, OpPop, OpLogicalOr, OpCopy, OpDup, OpImport:
		return 1
	case OpCall, OpReturn, OpStruct, OpMap, OpLabel, OpTuple, OpUntuple:
		return 2
	case OpAnnotation, OpPushConstant, OpPushLocalRef, OpPushForeign, OpStoreLocal, OpStoreForeign, OpArray:
		return 3
	case OpJmp, OpJnt:
		return 4
	default:
		panic("unknown op")
	}
}
func Sizes(ops ...Op) int {
	ret := 0
	for _, op := range ops {
		ret += Size(op)
	}
	return ret
}
func ReadFast(b []byte, p int, args []interface{}) (Op, int) {
	op := Op(b[p])
	switch op {
	case OpAdd, OpMult, OpGt, OpGte, OpLt, OpLte, OpClosure, OpSub, OpDiv, OpMod, OpEqTest, OpFieldAccess, OpFieldAssign, OpPop, OpLogicalOr, OpCopy, OpDup, OpImport:
		return op, 1
	case OpCall, OpReturn, OpStruct, OpMap, OpLabel, OpTuple, OpUntuple:
		args[0] = b[p+1]
		return op, 2
	case OpAnnotation, OpPushConstant, OpPushLocalRef, OpPushForeign, OpStoreLocal, OpStoreForeign, OpArray:
		args[0] = binary.BigEndian.Uint16(b[p+1:])
		return op, 3
	case OpJmp, OpJnt:
		args[0] = b[p+1]
		args[1] = binary.BigEndian.Uint16(b[p+2:])
		return op, 4
	default:
		panic("unknown op")
	}
}
func build(op Op, args args) (Instruction, error) {
	switch op {
	case OpAdd:
		return Add{}, nil
	case OpSub:
		return Sub{}, nil
	case OpGt:
		return Gt{}, nil
	case OpGte:
		return Gte{}, nil
	case OpLt:
		return Lt{}, nil
	case OpLte:
		return Lte{}, nil
	case OpPop:
		return Pop{}, nil
	case OpCopy:
		return Copy{}, nil
	case OpDup:
		return Dup{}, nil
	case OpImport:
		return Import{}, nil
	case OpLogicalOr:
		return LogicalOr{}, nil
	case OpPushConstant:
		id, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching constant index: %w", err)
		}
		return PushConstant{Index: id}, nil
	case OpPushLocalRef:
		id, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching local index: %w", err)
		}
		return PushLocal{Index: id}, nil
	case OpPushForeign:
		id, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching foreign index: %w", err)
		}
		return PushForeign{Index: id}, nil
	case OpJmp:
		addrType, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching addrtype: %w", err)
		}
		addr, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching addr: %w", err)
		}
		return Jmp{
			AddrType: addrType,
			Addr:     addr,
		}, nil
	case OpJnt:
		addrType, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching addrtype: %w", err)
		}
		addr, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching addr: %w", err)
		}
		return Jnt{
			AddrType: addrType,
			Addr:     addr,
		}, nil
	case OpAnnotation:
		id, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching constant index: %w", err)
		}
		return Annotation{
			Index: id,
		}, nil
	case OpMult:
		return Mult{}, nil
	case OpDiv:
		return Div{}, nil
	case OpMod:
		return Mod{}, nil
	case OpEqTest:
		return EqTest{}, nil
	case OpStoreLocal:
		id, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching local index: %w", err)
		}
		return StoreLocal{
			Index: id,
		}, nil
	case OpStoreForeign:
		id, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching foreign index: %w", err)
		}
		return StoreForeign{
			Index: id,
		}, nil
	case OpCall:
		argc, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching arg count: %w", err)
		}
		return Call{Args: argc}, nil
	case OpLabel:
		kind, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching label kind: %w", err)
		}
		return Label{Kind: LabelKind(kind)}, nil
	case OpArray:
		itemsc, err := args.Uint16()
		if err != nil {
			return nil, fmt.Errorf("fetching items count: %w", err)
		}
		return Array{Items: itemsc}, nil
	case OpTuple:
		itemsc, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching items count: %w", err)
		}
		return Tuple{Items: itemsc}, nil
	case OpUntuple:
		itemsc, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching items count: %w", err)
		}
		return Untuple{Items: itemsc}, nil
	case OpStruct:
		itemsc, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching items count: %w", err)
		}
		return Struct{Items: itemsc}, nil
	case OpMap:
		itemsc, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching items count: %w", err)
		}
		return Map{Items: itemsc}, nil
	case OpReturn:
		scope, err := args.Uint8()
		if err != nil {
			return nil, fmt.Errorf("fetching scope: %w", err)
		}
		return Return{Scope: scope}, nil
	case OpFieldAccess:
		return FieldAccess{}, nil
	case OpFieldAssign:
		return FieldAssign{}, nil
	case OpClosure:
		return Closure{}, nil
	default:
		return nil, fmt.Errorf("cannot build instruction for an unknown op: %s", op.String())
	}
}
func Bytes(i Instruction) []byte {
	if i == nil {
		return nil
	}
	switch inst := i.(type) {
	case Add, Gt, Lt, Gte, Lte, Mult, Closure, Sub, Div, Mod, EqTest, FieldAccess, FieldAssign, LogicalOr, Copy, Dup, Import:
		return bytes(inst.Op())
	case Call:
		return bytes(inst.Op(), inst.Args)
	case Label:
		return bytes(inst.Op(), uint8(inst.Kind))
	case Return:
		return bytes(inst.Op(), inst.Scope)
	case Array:
		return bytes(inst.Op(), inst.Items)
	case Tuple:
		return bytes(inst.Op(), inst.Items)
	case Untuple:
		return bytes(inst.Op(), inst.Items)
	case Struct:
		return bytes(inst.Op(), inst.Items)
	case Map:
		return bytes(inst.Op(), inst.Items)
	case PushConstant:
		return bytes(inst.Op(), inst.Index)
	case Pop:
		return bytes(inst.Op())
	case PushLocal:
		return bytes(inst.Op(), inst.Index)
	case PushForeign:
		return bytes(inst.Op(), inst.Index)
	case StoreLocal:
		return bytes(inst.Op(), inst.Index)
	case StoreForeign:
		return bytes(inst.Op(), inst.Index)
	case Jmp:
		return bytes(inst.Op(), inst.AddrType, inst.Addr)
	case Jnt:
		return bytes(inst.Op(), inst.AddrType, inst.Addr)
	case Annotation:
		return bytes(inst.Op(), inst.Index)
	default:
		panic("cannot generate code for an unknown instruction: " + reflect.TypeOf(i).String())
	}
}

func bytes(args ...interface{}) []byte {
	l := 0
	for _, a := range args {
		switch a.(type) {
		case Op:
			l += 1
		case uint8:
			l += 1
		case uint16:
			l += 2
		case uint32:
			l += 4
		default:
			panic("unsupported arg type: " + reflect.TypeOf(a).String())
		}
	}

	ret := make([]byte, l)
	offset := 0
	for _, a := range args {
		switch a := a.(type) {
		case Op:
			ret[offset] = byte(a)
			offset += 1
		case uint8:
			ret[offset] = a
			offset += 1
		case uint16:
			binary.BigEndian.PutUint16(ret[offset:], a)
			offset += 2
		case uint32:
			binary.BigEndian.PutUint32(ret[offset:], a)
			offset += 4
		default:
			panic("unsupported arg type")
		}
	}

	return ret
}
func Read(b []byte) (Instruction, int, error) {
	reader := &argsFromBytes{
		b:   b,
		pos: 0,
		cnt: 0,
	}
	if !reader.has(1) {
		return nil, 0, fmt.Errorf("cannot read op byte")
	}
	reader.pos++
	inst, err := build(Op(b[0]), reader)
	if err != nil {
		return nil, 0, fmt.Errorf("fetching arguments for %s: %w", Op(b[0]).String(), err)
	}

	return inst, reader.pos, nil
}
func Make(op Op, args ...interface{}) (Instruction, error) {
	reader := &argsFromArray{
		a: args,
	}
	inst, err := build(op, reader)
	if err != nil {
		return nil, fmt.Errorf("fetching arguments for %s: %w", op.String(), err)
	}

	if reader.pos < len(args) {
		return nil, fmt.Errorf("too many arguments for op %s", op.String())
	}

	return inst, nil
}
