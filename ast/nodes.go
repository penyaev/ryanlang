package ast

import (
	"fmt"
	"ryanlang/lexer"
	"strings"
)

type Expression interface {
	String() string
	Location() *lexer.Location
}

type Identifier struct {
	Name string
	Loc  *lexer.Location
}

func (i Identifier) Location() *lexer.Location {
	return i.Loc
}

func (i Identifier) String() string { return i.Name }

type String struct {
	Value string
	Loc   *lexer.Location
}

func (s String) Location() *lexer.Location {
	return s.Loc
}

func (s String) String() string { return "\"" + s.Value + "\"" }

type LetExpression struct {
	Identifiers    []Identifier
	Initialization Expression
	Loc            *lexer.Location
}

func (le LetExpression) Location() *lexer.Location {
	return le.Loc
}

func (le LetExpression) String() string {
	ids := []string{}
	for _, id := range le.Identifiers {
		ids = append(ids, id.String())
	}
	return fmt.Sprintf("let %s = %s", strings.Join(ids, ", "), le.Initialization.String())
}

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
}

func (ce CallExpression) Location() *lexer.Location {
	return ce.Callee.Location()
}

func (ce CallExpression) String() string {
	argumentStrings := []string{}
	for _, arg := range ce.Arguments {
		argumentStrings = append(argumentStrings, arg.String())
	}
	return fmt.Sprintf("%s(%s)", ce.Callee.String(), strings.Join(argumentStrings, ", "))
}

type NumberExpression struct {
	Value int
	Loc   *lexer.Location
}

func (n NumberExpression) Location() *lexer.Location {
	return n.Loc
}

func (n NumberExpression) String() string { return fmt.Sprintf("%v", n.Value) }

// todo: support multi-assign? e.g. `x, y = 1, 2`
type AssignExpression struct {
	Identifier Identifier
	Value      Expression
}

func (a AssignExpression) Location() *lexer.Location {
	return a.Identifier.Location()
}

func (a AssignExpression) String() string {
	return fmt.Sprintf("%s = %s", a.Identifier.String(), a.Value.String())
}

type FieldAssignExpression struct {
	FieldAccess FieldAccessExpression
	Value       Expression
}

func (f FieldAssignExpression) Location() *lexer.Location {
	return f.FieldAccess.Location()
}

func (f FieldAssignExpression) String() string {
	return fmt.Sprintf("%s = %s", f.FieldAccess.String(), f.Value.String())
}

type TupleAssignExpression struct {
	Tuple TupleExpression
	Value Expression
}

func (t TupleAssignExpression) Location() *lexer.Location {
	return t.Tuple.Location()
}

func (t TupleAssignExpression) String() string {
	return fmt.Sprintf("%s = %s", t.Tuple.String(), t.Value.String())
}

type IfExpression struct {
	Condition Expression
	Then      Expression
	Else      Expression
	Loc       *lexer.Location
}

func (i IfExpression) Location() *lexer.Location {
	return i.Loc
}

func (i IfExpression) String() string {
	var els string
	if i.Else != nil {
		els = fmt.Sprintf(" else %s", i.Else.String())
	}
	return fmt.Sprintf("if %s %s%s", i.Condition.String(), i.Then.String(), els)
}

type PlusExpression struct {
	Left  Expression
	Right Expression
}

func (p PlusExpression) Location() *lexer.Location {
	return p.Left.Location()
}

func (p PlusExpression) String() string {
	return fmt.Sprintf("(%s + %s)", p.Left.String(), p.Right.String())
}

type IncrExpression struct {
	Expr Expression
}

func (i IncrExpression) Location() *lexer.Location {
	return i.Expr.Location()
}

func (i IncrExpression) String() string {
	return fmt.Sprintf("(%s++)", i.Expr.String())
}

type MinusExpression struct {
	Left  Expression
	Right Expression
}

func (m MinusExpression) Location() *lexer.Location {
	return m.Left.Location()
}

func (m MinusExpression) String() string {
	return fmt.Sprintf("(%s - %s)", m.Left.String(), m.Right.String())
}

type DecrExpression struct {
	Expr Expression
}

func (d DecrExpression) Location() *lexer.Location {
	return d.Expr.Location()
}

func (d DecrExpression) String() string {
	return fmt.Sprintf("(%s--)", d.Expr.String())
}

type ModExpression struct {
	Left  Expression
	Right Expression
}

func (m ModExpression) Location() *lexer.Location {
	return m.Left.Location()
}

func (m ModExpression) String() string {
	return fmt.Sprintf("(%s %% %s)", m.Left.String(), m.Right.String())
}

type LogicalAndExpression struct {
	Left  Expression
	Right Expression
}

func (l LogicalAndExpression) Location() *lexer.Location {
	return l.Left.Location()
}

func (l LogicalAndExpression) String() string {
	return fmt.Sprintf("(%s && %s)", l.Left.String(), l.Right.String())
}

type LogicalOrExpression struct {
	Left  Expression
	Right Expression
}

func (l LogicalOrExpression) Location() *lexer.Location {
	return l.Left.Location()
}

func (l LogicalOrExpression) String() string {
	return fmt.Sprintf("(%s || %s)", l.Left.String(), l.Right.String())
}

type PrefixMinusExpression struct {
	Expr Expression
	loc  *lexer.Location
}

func (p PrefixMinusExpression) Location() *lexer.Location {
	return p.loc
}

func (p PrefixMinusExpression) String() string {
	return fmt.Sprintf("-%s", p.Expr.String())
}

type MultExpression struct {
	Left  Expression
	Right Expression
}

func (m MultExpression) Location() *lexer.Location {
	return m.Left.Location()
}

func (m MultExpression) String() string {
	return fmt.Sprintf("(%s * %s)", m.Left.String(), m.Right.String())
}

type DivExpression struct {
	Left  Expression
	Right Expression
}

func (d DivExpression) Location() *lexer.Location {
	return d.Left.Location()
}

func (d DivExpression) String() string {
	return fmt.Sprintf("(%s / %s)", d.Left.String(), d.Right.String())
}

type GtExpression struct {
	Left  Expression
	Right Expression
}

func (g GtExpression) Location() *lexer.Location {
	return g.Left.Location()
}

func (g GtExpression) String() string {
	return fmt.Sprintf("(%s > %s)", g.Left.String(), g.Right.String())
}

type LtExpression struct {
	Left  Expression
	Right Expression
}

func (l LtExpression) Location() *lexer.Location {
	return l.Left.Location()
}

func (l LtExpression) String() string {
	return fmt.Sprintf("(%s < %s)", l.Left.String(), l.Right.String())
}

type GteExpression struct {
	Left  Expression
	Right Expression
}

func (g GteExpression) Location() *lexer.Location {
	return g.Left.Location()
}

func (g GteExpression) String() string {
	return fmt.Sprintf("(%s >= %s)", g.Left.String(), g.Right.String())
}

type LteExpression struct {
	Left  Expression
	Right Expression
}

func (l LteExpression) Location() *lexer.Location {
	return l.Left.Location()
}

func (l LteExpression) String() string {
	return fmt.Sprintf("(%s <= %s)", l.Left.String(), l.Right.String())
}

type TupleExpression struct {
	Exprs []Expression
}

func (t TupleExpression) Location() *lexer.Location {
	return t.Exprs[0].Location()
}

func (t TupleExpression) String() string {
	s := []string{}
	for _, expr := range t.Exprs {
		s = append(s, expr.String())
	}
	return fmt.Sprintf("tuple(%s)", strings.Join(s, ", "))
}

type EqTestExpression struct {
	Left  Expression
	Right Expression
}

func (e EqTestExpression) Location() *lexer.Location {
	return e.Left.Location()
}

func (e EqTestExpression) String() string {
	return fmt.Sprintf("(%s == %s)", e.Left.String(), e.Right.String())
}

type FuncExpression struct {
	Arguments []Identifier
	Body      Expression
	Loc       *lexer.Location
}

func (f FuncExpression) Location() *lexer.Location {
	return f.Loc
}

func (f FuncExpression) String() string {
	argumentsStrings := []string{}
	for _, arg := range f.Arguments {
		argumentsStrings = append(argumentsStrings, arg.String())
	}

	return fmt.Sprintf("func(%s) %s", strings.Join(argumentsStrings, ", "), f.Body.String())
}

type ReturnExpression struct {
	Expr Expression
	loc  *lexer.Location
}

func (r ReturnExpression) Location() *lexer.Location {
	return r.loc
}

func (r ReturnExpression) String() string {
	return fmt.Sprintf("return %s", r.Expr.String())
}

type ContinueExpression struct {
	loc *lexer.Location
}

func (c ContinueExpression) Location() *lexer.Location {
	return c.loc
}

func (c ContinueExpression) String() string {
	return "continue"
}

type BreakExpression struct {
	loc *lexer.Location
}

func (b BreakExpression) Location() *lexer.Location {
	return b.loc
}

func (b BreakExpression) String() string {
	return "break"
}

type BlockExpression struct {
	Stmts []Statement
	Loc   *lexer.Location
}

func (b BlockExpression) Location() *lexer.Location {
	return b.Loc
}

func (b BlockExpression) String() string {
	s := []string{}
	for _, st := range b.Stmts {
		s = append(s, st.String())
	}
	return "{ " + strings.Join(s, " ") + " }"
}

type GroupExpression struct {
	Expr Expression
	Loc  *lexer.Location
}

func (g GroupExpression) Location() *lexer.Location {
	return g.Loc
}

func (g GroupExpression) String() string {
	return "(" + g.Expr.String() + ")"
}

type Statement struct {
	Expr Expression
}

func (s Statement) Location() *lexer.Location {
	return s.Expr.Location()
}

func (s Statement) String() string {
	return s.Expr.String() + ";"
}

type BuiltinFunction struct {
	Name string
	loc  *lexer.Location
}

func (b BuiltinFunction) Location() *lexer.Location {
	return b.loc
}

func (b BuiltinFunction) String() string {
	return fmt.Sprintf("(built-in function %s)", b.Name)
}

type StructExpression struct {
	Fields map[string]Expression
	loc    *lexer.Location
}

func (s StructExpression) Location() *lexer.Location {
	return s.loc
}

func (s StructExpression) String() string {
	strs := []string{}
	for key, value := range s.Fields {
		strs = append(strs, fmt.Sprintf("%s: %s", key, value.String()))
	}
	return fmt.Sprintf("struct { %s }", strings.Join(strs, "; "))
}

type MapField struct {
	Key   Expression
	Value Expression
}
type MapExpression struct {
	Fields []MapField
	loc    *lexer.Location
}

func (m MapExpression) Location() *lexer.Location {
	return m.loc
}

func (m MapExpression) String() string {
	strs := []string{}
	for _, v := range m.Fields {
		strs = append(strs, fmt.Sprintf("%s: %s", v.Key.String(), v.Value.String()))
	}
	return fmt.Sprintf("map { %s }", strings.Join(strs, "; "))
}

type FieldAccessExpression struct {
	Left  Expression
	Right Expression
}

func (f FieldAccessExpression) Location() *lexer.Location {
	return f.Left.Location()
}

func (f FieldAccessExpression) String() string {
	return fmt.Sprintf("%s.%s", f.Left.String(), f.Right.String())
}

// todo: support fors?
// todo: support break/continue
type WhileExpression struct {
	Condition Expression
	Body      Expression
	Loc       *lexer.Location
}

func (w WhileExpression) Location() *lexer.Location {
	return w.Loc
}

func (w WhileExpression) String() string {
	return fmt.Sprintf("while %s %s", w.Condition.String(), w.Body.String())
}

type ForExpression struct {
	Index *Identifier
	Value Identifier
	Range Expression
	Body  Expression
	Loc   *lexer.Location
}

func (f ForExpression) Location() *lexer.Location {
	return f.Loc
}

func (f ForExpression) String() string {
	iterators := []string{}
	if f.Index != nil {
		iterators = append(iterators, f.Index.String())
	}
	iterators = append(iterators, f.Value.String())
	return fmt.Sprintf("for %s in %s %s", strings.Join(iterators, ", "), f.Range.String(), f.Body.String())
}

// todo: support hashes (or somehow combine them with structs?)
type ArrayExpression struct {
	Items []Expression
	loc   *lexer.Location
}

func (a ArrayExpression) Location() *lexer.Location {
	return a.loc
}

func (a ArrayExpression) String() string {
	argumentStrings := []string{}
	for _, item := range a.Items {
		argumentStrings = append(argumentStrings, item.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(argumentStrings, ", "))
}

type NegationExpression struct {
	Expr Expression
	loc  *lexer.Location
}

func (n NegationExpression) Location() *lexer.Location {
	return n.loc
}

func (n NegationExpression) String() string {
	return fmt.Sprintf("!%s", n.Expr.String())
}

type Module struct {
	Name  string
	Block BlockExpression
}

func (m Module) Location() *lexer.Location {
	return m.Block.Location()
}

func (m Module) String() string {
	return fmt.Sprintf("(module %s)", m.Name)
}

type Exports struct {
	Fields map[string]Expression
	loc    *lexer.Location
}

func (e Exports) Location() *lexer.Location {
	return e.loc
}

func (e Exports) String() string {
	fields := []string{}
	for k := range e.Fields {
		fields = append(fields, k)
	}
	return fmt.Sprintf("(exports %s)", strings.Join(fields, ", "))
}

type Import struct {
	Module Expression
	Loc    *lexer.Location
}

func (i Import) Location() *lexer.Location {
	return i.Loc
}

func (i Import) String() string {
	return fmt.Sprintf("import(%s)", i.Module.String())
}

type ArrowExpression struct {
	Expr Expression
	Loc  *lexer.Location
}

func (a ArrowExpression) Location() *lexer.Location {
	return a.Loc
}

func (a ArrowExpression) String() string {
	return fmt.Sprintf("=> %s", a.Expr.String())
}
