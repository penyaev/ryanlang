package instruction

import "strconv"

type Op byte

const (
	OpUnknown Op = iota
	OpAdd
	OpSub
	OpMult
	OpDiv
	OpMod
	OpGt
	OpGte
	OpLt
	OpLte
	OpLogicalOr
	OpEqTest
	OpPushConstant // todo: allow pushing values fitting into uint16 directly, w/out using the const pool, e.g. PushNumber
	// todo: maybe introduce ops like push8, push16, push32 for different operand sizes?
	OpPushLocalRef
	OpCopy
	OpDup
	OpPushForeign
	OpPop
	OpClosure
	OpJnt
	OpJmp
	OpAnnotation
	OpStoreLocal
	OpStoreForeign
	OpCall
	OpReturn
	OpArray
	OpTuple
	OpUntuple
	OpStruct
	OpMap
	OpFieldAccess
	OpFieldAssign
	OpImport
	OpLabel
)

func (o Op) String() string {
	switch o {
	case OpUnknown:
		return "???"
	case OpAdd:
		return "ADD"
	case OpSub:
		return "SUB"
	case OpMult:
		return "MULT"
	case OpDiv:
		return "DIV"
	case OpMod:
		return "MOD"
	case OpGt:
		return "GT"
	case OpGte:
		return "GTE"
	case OpLt:
		return "LT"
	case OpLogicalOr:
		return "LOR"
	case OpLte:
		return "LTE"
	case OpEqTest:
		return "EQ"
	case OpPushConstant:
		return "PUSHC"
	case OpPushLocalRef:
		return "PUSHL"
	case OpCopy:
		return "COPY"
	case OpDup:
		return "DUP"
	case OpPushForeign:
		return "PUSHF"
	case OpPop:
		return "POP"
	case OpJnt:
		return "JNT"
	case OpJmp:
		return "JMP"
	case OpStoreLocal:
		return "STOREL"
	case OpStoreForeign:
		return "STOREF"
	case OpAnnotation:
		return "TXT"
	case OpCall:
		return "CALL"
	case OpClosure:
		return "CLOSURE"
	case OpReturn:
		return "RETURN"
	case OpArray:
		return "ARRAY"
	case OpTuple:
		return "TUPLE"
	case OpUntuple:
		return "UNTUPLE"
	case OpStruct:
		return "STRUCT"
	case OpMap:
		return "MAP"
	case OpFieldAccess:
		return "PUSHFLD"
	case OpFieldAssign:
		return "STOREFLD"
	case OpImport:
		return "IMPORT"
	case OpLabel:
		return "LABEL"
	default:
		panic("cannot stringify unknown op: " + strconv.Itoa(int(o)))
	}
}
