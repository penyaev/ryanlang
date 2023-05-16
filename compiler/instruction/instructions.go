package instruction

import (
	"fmt"
)

type Instruction interface {
	Op() Op
	String() string
}

type Annotation struct {
	Index uint16
}

func (Annotation) Op() Op {
	return OpAnnotation
}
func (a Annotation) String() string {
	return fmt.Sprintf("%s\t%d", a.Op().String(), a.Index)
}

type Mult struct{}

func (Mult) Op() Op {
	return OpMult
}
func (m Mult) String() string {
	return fmt.Sprintf("%s", m.Op().String())
}

type LogicalOr struct{}

func (LogicalOr) Op() Op {
	return OpLogicalOr
}
func (l LogicalOr) String() string {
	return fmt.Sprintf("%s", l.Op().String())
}

type EqTest struct{}

func (EqTest) Op() Op {
	return OpEqTest
}
func (e EqTest) String() string {
	return fmt.Sprintf("%s", e.Op().String())
}

type Div struct{}

func (Div) Op() Op {
	return OpDiv
}
func (d Div) String() string {
	return fmt.Sprintf("%s", d.Op().String())
}

type Mod struct{}

func (Mod) Op() Op {
	return OpMod
}
func (m Mod) String() string {
	return fmt.Sprintf("%s", m.Op().String())
}

type Gt struct{}

func (Gt) Op() Op {
	return OpGt
}
func (g Gt) String() string {
	return fmt.Sprintf("%s", g.Op().String())
}

type Gte struct{}

func (Gte) Op() Op {
	return OpGte
}
func (g Gte) String() string {
	return fmt.Sprintf("%s", g.Op().String())
}

type Lt struct{}

func (Lt) Op() Op {
	return OpLt
}
func (l Lt) String() string {
	return fmt.Sprintf("%s", l.Op().String())
}

type Lte struct{}

func (Lte) Op() Op {
	return OpLte
}
func (l Lte) String() string {
	return fmt.Sprintf("%s", l.Op().String())
}

type Jnt struct {
	AddrType uint8
	Addr     uint16
}

func (Jnt) Op() Op {
	return OpJnt
}
func (j Jnt) String() string {
	return fmt.Sprintf("%s\t%d %d", j.Op().String(), j.AddrType, j.Addr)
}

type Jmp struct {
	AddrType uint8
	Addr     uint16
}

func (Jmp) Op() Op {
	return OpJmp
}
func (j Jmp) String() string {
	return fmt.Sprintf("%s\t%d %d", j.Op().String(), j.AddrType, j.Addr)
}

type PushConstant struct {
	Index uint16
}

func (PushConstant) Op() Op {
	return OpPushConstant
}
func (p PushConstant) String() string {
	return fmt.Sprintf("%s\t%d", p.Op().String(), p.Index)
}

type PushLocal struct {
	Index uint16
}

func (PushLocal) Op() Op {
	return OpPushLocalRef
}
func (p PushLocal) String() string {
	return fmt.Sprintf("%s\t%d", p.Op().String(), p.Index)
}

type PushForeign struct {
	Index uint16
}

func (PushForeign) Op() Op {
	return OpPushForeign
}
func (p PushForeign) String() string {
	return fmt.Sprintf("%s\t%d", p.Op().String(), p.Index)
}

type Pop struct {
}

func (Pop) Op() Op {
	return OpPop
}
func (p Pop) String() string {
	return fmt.Sprintf("%s", p.Op().String())
}

type Add struct{}

func (a Add) Op() Op {
	return OpAdd
}
func (a Add) String() string {
	return fmt.Sprintf("%s", a.Op().String())
}

type Sub struct{}

func (s Sub) Op() Op {
	return OpSub
}
func (s Sub) String() string {
	return fmt.Sprintf("%s", s.Op().String())
}

type Return struct {
	Scope uint8
}

func (r Return) Op() Op {
	return OpReturn
}
func (r Return) String() string {
	return fmt.Sprintf("%s\t%d", r.Op().String(), r.Scope)
}

type Closure struct{}

func (c Closure) Op() Op {
	return OpClosure
}
func (c Closure) String() string {
	return fmt.Sprintf("%s", c.Op().String())
}

type Copy struct{}

func (c Copy) Op() Op {
	return OpCopy
}
func (c Copy) String() string {
	return fmt.Sprintf("%s", c.Op().String())
}

type Dup struct{}

func (d Dup) Op() Op {
	return OpDup
}
func (d Dup) String() string {
	return fmt.Sprintf("%s", d.Op().String())
}

type FieldAccess struct{}

func (f FieldAccess) Op() Op {
	return OpFieldAccess
}
func (f FieldAccess) String() string {
	return fmt.Sprintf("%s", f.Op().String())
}

type FieldAssign struct{}

func (f FieldAssign) Op() Op {
	return OpFieldAssign
}
func (f FieldAssign) String() string {
	return fmt.Sprintf("%s", f.Op().String())
}

type StoreLocal struct {
	Index uint16
}

func (StoreLocal) Op() Op {
	return OpStoreLocal
}
func (s StoreLocal) String() string {
	return fmt.Sprintf("%s\t%d", s.Op().String(), s.Index)
}

type StoreForeign struct {
	Index uint16
}

func (StoreForeign) Op() Op {
	return OpStoreForeign
}
func (s StoreForeign) String() string {
	return fmt.Sprintf("%s\t%d", s.Op().String(), s.Index)
}

type Call struct {
	Args uint8
}

func (Call) Op() Op {
	return OpCall
}
func (c Call) String() string {
	return fmt.Sprintf("%s\t%d", c.Op().String(), c.Args)
}

type Array struct {
	Items uint16
}

func (Array) Op() Op {
	return OpArray
}
func (a Array) String() string {
	return fmt.Sprintf("%s\t%d", a.Op().String(), a.Items)
}

type Tuple struct {
	Items uint8 // when changing datatype here, also change it for Untuple op
}

func (Tuple) Op() Op {
	return OpTuple
}
func (t Tuple) String() string {
	return fmt.Sprintf("%s\t%d", t.Op().String(), t.Items)
}

type Untuple struct {
	Items uint8 // when changing datatype here, also change it for Tuple op
}

func (Untuple) Op() Op {
	return OpUntuple
}
func (u Untuple) String() string {
	return fmt.Sprintf("%s\t%d", u.Op().String(), u.Items)
}

type LabelKind uint8

const (
	LabelKindContinue LabelKind = iota
)

type Label struct {
	Kind LabelKind
}

func (Label) Op() Op {
	return OpLabel
}
func (l Label) String() string {
	return fmt.Sprintf("%s\t%d", l.Op().String(), l.Kind)
}

type Struct struct {
	Items uint8
}

func (Struct) Op() Op {
	return OpStruct
}
func (s Struct) String() string {
	return fmt.Sprintf("%s\t%d", s.Op().String(), s.Items)
}

type Map struct {
	Items uint8
}

func (Map) Op() Op {
	return OpMap
}
func (m Map) String() string {
	return fmt.Sprintf("%s\t%d", m.Op().String(), m.Items)
}

type Import struct {
}

func (Import) Op() Op {
	return OpImport
}
func (i Import) String() string {
	return fmt.Sprintf("%s", i.Op().String())
}
