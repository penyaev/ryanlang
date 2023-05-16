package compiler

import (
	"fmt"
	"ryanlang/ast"
	"ryanlang/compiler/instruction"
	"ryanlang/compiler/sourcemap"
	"ryanlang/funcs"
	"ryanlang/lexer"
	"ryanlang/object"
)

const (
	RelativeAddress uint8 = iota
	AbsoluteAddress
	RelativeAddressBackwards
	LabelAddress
)

type DebugSymbolData struct {
	Name string
	Loc  *lexer.Location
}
type DebugData struct {
	sm       *sourcemap.SourceMap
	Locals   []DebugSymbolData
	Foreigns []DebugSymbolData
}

func (dd *DebugData) SearchSource(p int) *sourcemap.Entry {
	if dd == nil || dd.sm == nil {
		return nil
	}
	return dd.sm.Search(p)
}
func (dd *DebugData) LocalName(id int) string {
	return dd.Locals[id].Name
}
func (dd *DebugData) ForeignName(id int) string {
	return dd.Foreigns[id].Name
}

type Compiler struct {
	objects   *object.Storage
	scopes    []*scope
	symbols   *symbols
	debugData map[int]*DebugData
}

func NewCompilerWithStorage(objects *object.Storage) *Compiler {
	c := &Compiler{
		scopes:    []*scope{newScope()},
		symbols:   newSymbols(),
		objects:   objects,
		debugData: map[int]*DebugData{},
	}
	c.registerBuiltins()
	return c
}

func NewCompiler() *Compiler {
	return NewCompilerWithStorage(object.NewStorage())
}
func (c *Compiler) saveDebugData(objId int, sm *sourcemap.SourceMap, sym *symbols) {
	localsData := make([]DebugSymbolData, *sym.locals)
	foreignsData := make([]DebugSymbolData, len(sym.foreign))

	for _, s := range sym.storage {
		if s.scope == symbolScopeLocal {
			localsData[s.id] = DebugSymbolData{
				Name: s.name,
				Loc:  nil, // todo: add location
			}
		}
		if s.scope == symbolScopeForeign {
			foreignsData[s.id] = DebugSymbolData{
				Name: s.name,
				Loc:  nil, // todo: add location
			}
		}
	}
	c.debugData[objId] = &DebugData{
		sm:       sm,
		Locals:   localsData,
		Foreigns: foreignsData,
	}
}

func (c *Compiler) CompileRunModule(node ast.Module) (*Module, error) {
	//c.enterScope()
	module, err := c.CompileModule(node)
	if err != nil {
		return nil, err
	}
	err = iferr(
		err,
		c.annotate("entry point for module: "+node.Name),
		c.emitInstruction(instruction.OpPushConstant, int(module.EntryPoint)),
		c.emitInstruction(instruction.OpClosure),
		c.emitInstruction(instruction.OpCall, 0),
		c.emitInstruction(instruction.OpReturn, int(object.CodeReturnScopeFunc)),
	)
	//s := c.leaveScope()
	if err != nil {
		return nil, err
	}
	id := c.registerObject(&object.Code{
		Code:        c.scope().code.b,
		Locals:      *c.symbols.locals,
		Foreigns:    0,
		Loc:         node.Location(), // todo: debug only?
		ReturnScope: object.CodeReturnScopeFunc,
	})
	c.saveDebugData(id, c.scope().code.sm, c.symbols)

	compiledModule := &Module{
		EntryPoint: uint16(id),
		Objects:    c.objects,
		DebugData:  c.debugData,
	}
	return compiledModule, nil
}

func (c *Compiler) annotate(s string) error {
	return c.emitInstruction(instruction.OpAnnotation, c.registerObject(&object.String{Value: s}))
}
func (c *Compiler) scope() *scope {
	return c.scopes[len(c.scopes)-1]
}
func (c *Compiler) enterScope() *scope {
	c.scopes = append(c.scopes, newScope())
	return c.scope()
}
func (c *Compiler) leaveScope() *scope {
	cur := c.scope()
	cur.code.sm.PopAll(c.pos())
	c.scopes[len(c.scopes)-1] = nil
	c.scopes = c.scopes[:len(c.scopes)-1]
	return cur
}

func (c *Compiler) pushSymbols() {
	c.symbols = c.symbols.push()
}
func (c *Compiler) pushSymbolsLinked() {
	c.symbols = c.symbols.pushLinked()
}
func (c *Compiler) popSymbols() *symbols {
	prev := c.symbols
	c.symbols = c.symbols.pop()
	return prev
}

func (c *Compiler) registerObject(obj object.Object) int {

	return int(c.objects.Add(obj))
}
func (c *Compiler) emitConstantObject(obj object.Object) error {
	return c.emitInstruction(instruction.OpPushConstant, c.registerObject(obj))
}

func (c *Compiler) emitPushBuiltin(key string) error {
	s := c.symbols.getWithScope(key, symbolScopeBuiltin)
	if s == nil {
		panic("cannot get built-in " + key)
	}
	return c.emitInstruction(instruction.OpPushConstant, s.id)
}
func (c *Compiler) emitPushNull() error {
	return c.emitPushBuiltin("null")
}
func (c *Compiler) emitPushTrue() error {
	return c.emitPushBuiltin("true")
}
func (c *Compiler) emitPushFalse() error {
	return c.emitPushBuiltin("false")
}
func (c *Compiler) makePushBuiltin(key string) (*code, error) {
	return c.makecb(func() error {
		return c.emitPushBuiltin(key)
	})
}
func (c *Compiler) makePushNull() (*code, error) {
	return c.makePushBuiltin("null")
}
func (c *Compiler) makePushTrue() (*code, error) {
	return c.makePushBuiltin("true")
}
func (c *Compiler) makePushFalse() (*code, error) {
	return c.makePushBuiltin("false")
}

func (c *Compiler) makeInstruction(op instruction.Op, args ...interface{}) (*code, error) {
	inst, err := instruction.Make(op, args...)
	if err != nil {
		return nil, err
	}
	b := instruction.Bytes(inst)
	if len(b) != instruction.Size(op) {
		return nil, fmt.Errorf("instruction %s size is expected to be %d, got %d", op.String(), instruction.Size(op), len(b))
	}
	return &code{b: b}, nil
}
func (c *Compiler) emitInstruction(op instruction.Op, args ...interface{}) error {
	b, err := c.makeInstruction(op, args...)
	if err != nil {
		return err
	}
	return c.emit(b)
}

func (c *Compiler) makeClosure(maker func() (*code, error), args []ast.Identifier, rs object.CodeReturnScope) (*code, error) {
	c.pushSymbols()
	for _, arg := range args {
		c.symbols.createLocal(arg.Name)
	}
	bodyCode, err := maker()
	sym := c.popSymbols()
	if err != nil {
		return nil, err
	}

	id := c.registerObject(&object.Code{
		Code:        bodyCode.b,
		Locals:      *sym.locals,
		Arguments:   len(args),
		Foreigns:    len(sym.foreign),
		ReturnScope: rs,
	})
	c.saveDebugData(id, bodyCode.sm, sym)

	err = nil
	var closureCode *code
	closureCode, err = c.makecb(func() error {
		for i := len(sym.foreign) - 1; i >= 0; i-- {
			err = iferr(err, c.emitPushSymbol(sym.foreign[i]))
		}

		err = iferr(
			err,
			c.emitInstruction(instruction.OpPushConstant, id),
			c.emitInstruction(instruction.OpClosure),
		)
		return err
	})
	return closureCode, err
}
func (c *Compiler) emitPushSymbol(s *Symbol) error {
	var err error
	switch s.scope {
	case symbolScopeBuiltin:
		err = iferr(err, c.emitInstruction(instruction.OpPushConstant, s.id))
	case symbolScopeLocal:
		err = iferr(err, c.emitInstruction(instruction.OpPushLocalRef, s.id))
	case symbolScopeForeign:
		err = iferr(err, c.emitInstruction(instruction.OpPushForeign, s.id))
	}
	return err
}
func (c *Compiler) emitStoreSymbol(s *Symbol) error {
	var err error
	switch s.scope {
	case symbolScopeLocal:
		err = iferr(err, c.emitInstruction(instruction.OpStoreLocal, s.id))
	case symbolScopeForeign:
		err = iferr(err, c.emitInstruction(instruction.OpStoreForeign, s.id))
	default:
		return fmt.Errorf("cannot store to symbol with scope %s", s.scope.String())
	}
	return err
}
func (c *Compiler) emit(args ...*code) error {
	for _, a := range args {
		if a == nil {
			continue
		}
		if a.sm != nil && c.scope().code.sm != nil {
			c.scope().code.sm.Merge(a.sm, c.pos())
		}
		c.scope().code.b = append(c.scope().code.b, a.b...)
	}

	return nil
}
func (c *Compiler) pos() int {
	return len(c.scope().code.b)
}
func (c *Compiler) makecb(emitter func() error) (*code, error) {
	c.enterScope()
	err := emitter()
	s := c.leaveScope()
	if err != nil {
		return nil, err
	}
	return s.code, nil
}
func (c *Compiler) make(node ast.Expression) (*code, error) {
	return c.makecb(func() error {
		return c.emitNode(node)
	})
}
func (c *Compiler) makeExpressions(nodes []ast.Expression) (*code, error) {
	return c.makecb(func() error {
		return c.compileExpressions(nodes)
	})
}
func (c *Compiler) makeStatements(nodes []ast.Statement) (*code, error) {
	return c.makecb(func() error {
		return c.compileStatements(nodes)
	})
}
func (c *Compiler) registerBuiltins() {
	for key, obj := range funcs.Builtins {
		c.symbols.createBuiltin(key, c.registerObject(obj))
	}
	for key := range funcs.BuiltinFunctions {
		c.symbols.createBuiltin(key, c.registerObject(&object.Closure{BuiltinFunctionName: key}))
	}
}
