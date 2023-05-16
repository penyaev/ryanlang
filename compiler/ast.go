package compiler

import (
	"fmt"
	"reflect"
	"ryanlang/ast"
	"ryanlang/compiler/instruction"
	"ryanlang/object"
)

func (c *Compiler) compileNumberExpression(node ast.NumberExpression) error {
	return c.emitConstantObject(&object.Number{Value: node.Value})
}
func (c *Compiler) compileArrayExpression(node ast.ArrayExpression) error {
	return iferr(
		c.compileExpressionsReversed(node.Items),
		c.emitInstruction(instruction.OpArray, len(node.Items)),
	)
}
func (c *Compiler) compileTupleExpression(node ast.TupleExpression) error {
	return iferr(
		c.compileExpressionsReversed(node.Exprs),
		c.emitInstruction(instruction.OpTuple, len(node.Exprs)),
	)
}
func (c *Compiler) compileStructExpression(node ast.StructExpression) error {
	c.pushSymbolsLinked()
	var th *Symbol
	var err error
	for key, expr := range node.Fields {
		isFunction := false
		if _, ok := expr.(ast.FuncExpression); ok {
			c.pushSymbolsLinked()
			isFunction = true
			if th == nil {
				th = c.symbols.createLocal("this")
			} else {
				c.symbols.importLocal(th)
			}
		}
		err = iferr(
			err,
			c.emitNode(expr),
			c.emitConstantObject(&object.String{Value: key}),
		)

		if isFunction {
			c.popSymbols()
		}
	}
	c.popSymbols()
	err = iferr(
		err,
		c.emitInstruction(instruction.OpStruct, len(node.Fields)),
		//c.emitStoreSymbol(th),
		//c.emitPushSymbol(th),
	)
	if th != nil {
		err = iferr(err, c.emitStoreSymbol(th))
	}
	return err
}
func (c *Compiler) compileMapExpression(node ast.MapExpression) error {
	c.pushSymbolsLinked()
	var err error
	for _, field := range node.Fields {
		err = iferr(
			err,
			c.emitNode(field.Value),
			c.emitNode(field.Key),
		)
	}
	c.popSymbols()
	return iferr(
		err,
		c.emitInstruction(instruction.OpMap, len(node.Fields)),
	)
}
func (c *Compiler) compileFieldAccessExpression(node ast.FieldAccessExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpFieldAccess),
	)
}
func (c *Compiler) compileFieldAssignExpression(node ast.FieldAssignExpression) error {
	return iferr(
		c.emitNode(node.Value),
		c.emitNode(node.FieldAccess.Right),
		c.emitNode(node.FieldAccess.Left),
		c.emitInstruction(instruction.OpFieldAssign),
	)
}
func (c *Compiler) compileString(node ast.String) error {
	return c.emitConstantObject(&object.String{Value: node.Value})
}
func (c *Compiler) compilePlusExpression(node ast.PlusExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpAdd),
	)
}
func (c *Compiler) compileMinusExpression(node ast.MinusExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpSub),
	)
}
func (c *Compiler) compileMultExpression(node ast.MultExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpMult),
	)
}
func (c *Compiler) compileDivExpression(node ast.DivExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpDiv),
	)
}
func (c *Compiler) compileModExpression(node ast.ModExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpMod),
	)
}
func (c *Compiler) compileGtExpression(node ast.GtExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpGt),
	)
}
func (c *Compiler) compileGteExpression(node ast.GteExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpGte),
	)
}
func (c *Compiler) compileLtExpression(node ast.LtExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpLt),
	)
}
func (c *Compiler) compileLteExpression(node ast.LteExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpLte),
	)
}
func (c *Compiler) compileLogicalAndExpression(node ast.LogicalAndExpression) error {
	// key point is to first eval left, and only eval right if left is true
	jnt, err1 := c.makeInstruction(instruction.OpJnt, int(RelativeAddress), 0) // todo: cleaner
	jmp, err2 := c.makeInstruction(instruction.OpJmp, int(RelativeAddress), 0)
	pt, err3 := c.makePushTrue()
	pf, err4 := c.makePushFalse()
	left, err5 := c.make(node.Left)
	right, err6 := c.make(node.Right)
	return iferr(
		err1, err2, err3, err4, err5, err6,
		c.emit(left),
		c.emitInstruction(instruction.OpJnt, int(RelativeAddress), right.Len()+jnt.Len()+pt.Len()+jmp.Len()),
		c.emit(right),
		c.emitInstruction(instruction.OpJnt, int(RelativeAddress), pt.Len()+jmp.Len()),
		c.emit(pt),
		c.emitInstruction(instruction.OpJmp, int(RelativeAddress), pf.Len()),
		c.emit(pf),
	)
}
func (c *Compiler) compileLogicalOrExpression(node ast.LogicalOrExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpLogicalOr),
	)
}
func (c *Compiler) compileEqTestExpression(node ast.EqTestExpression) error {
	return iferr(
		c.emitNode(node.Right),
		c.emitNode(node.Left),
		c.emitInstruction(instruction.OpEqTest),
	)
}
func (c *Compiler) compileExpressions(node []ast.Expression) error {
	for _, expr := range node {
		if err := c.emitNode(expr); err != nil {
			return err
		}
	}
	return nil
}
func (c *Compiler) compileCopiedExpressions(node []ast.Expression) error {
	for _, expr := range node {
		if err := iferr(c.emitNode(expr), c.emitInstruction(instruction.OpCopy)); err != nil {
			return err
		}
	}
	return nil
}
func (c *Compiler) compileStatements(node []ast.Statement) error {
	for _, expr := range node {
		if err := c.emitNode(expr); err != nil {
			return err
		}

	}
	return nil
}
func (c *Compiler) compileExpressionsReversed(node []ast.Expression) error {
	for i := len(node) - 1; i >= 0; i-- {
		if err := c.emitNode(node[i]); err != nil {
			return err
		}
	}
	return nil
}
func (c *Compiler) compileIfExpression(node ast.IfExpression) error {
	c.pushSymbolsLinked()
	condition, err := c.make(node.Condition)
	c.pushSymbolsLinked()
	conseq, err1 := c.make(node.Then)
	c.popSymbols()

	var alt, jmp *code
	var err2, err3 error
	if node.Else != nil {
		c.pushSymbolsLinked()
		alt, err2 = c.make(node.Else)
		c.popSymbols()
	} else {
		alt, err2 = c.makePushNull() // if there's no "else" block, If expression still should resolve to something
	}
	c.popSymbols()

	jmp, err3 = c.makeInstruction(instruction.OpJmp, int(RelativeAddress), alt.Len())
	return iferr(
		err,
		err1,
		err2,
		err3,
		c.emit(condition), // todo: should be in nested scope
		c.emitInstruction(instruction.OpJnt, int(RelativeAddress), conseq.Len()+jmp.Len()),
		c.emit(conseq),
		c.emit(jmp),
		c.emit(alt),
	)
}
func (c *Compiler) compileWhileExpressionBody(node ast.WhileExpression) error {
	c.pushSymbolsLinked()
	condition, err := c.make(node.Condition)
	c.pushSymbolsLinked()
	body, err2 := c.make(node.Body)
	c.popSymbols()
	c.popSymbols()

	err = iferr(
		err, err2,
		c.annotate("inner while start"),
		c.emitPushNull(), // if there'll be no iterations of the loop, it will resolve to this null
	)
	//loopStart := c.pos()
	err = iferr(
		err,
		c.emitInstruction(instruction.OpLabel, int(instruction.LabelKindContinue)), // "loop start"
		c.emit(condition),
		c.emitInstruction(instruction.OpJnt, int(RelativeAddress), instruction.Size(instruction.OpPop)+body.Len()+instruction.Size(instruction.OpJmp)),
		c.emitInstruction(instruction.OpPop), // remove the resolve value, because we'll now get another one by executing body
		c.emit(body),
		// push label "continue"
		// jmp to "loop start"
		//c.emitInstruction(instruction.OpJmp, RelativeAddressBackwards, (c.dist(loopStart)+instruction.Size(instruction.OpJmp))),
		c.emitInstruction(instruction.OpJmp, int(LabelAddress), int(instruction.LabelKindContinue)),
		// push label "break"
		// pop label "break"
		// pop label "continue"
		// pop label "loop start"
	)
	return err
}
func (c *Compiler) compileNegationExpression(node ast.NegationExpression) error {
	return iferr(
		c.emitNode(node.Expr),
		c.emitInstruction(instruction.OpJnt, int(RelativeAddress), instruction.Sizes(instruction.OpPushConstant, instruction.OpJmp)),
		//c.emitInstruction(instruction.OpPop),
		c.emitPushFalse(),
		c.emitInstruction(instruction.OpJmp, int(RelativeAddress), instruction.Sizes(instruction.OpPushConstant)),
		//c.emitInstruction(instruction.OpPop),
		c.emitPushTrue(),
	)
}
func (c *Compiler) compilePrefixMinusExpression(node ast.PrefixMinusExpression) error {
	return iferr(
		c.emitNode(node.Expr),
		c.emitConstantObject(&object.Number{Value: 0}),
		c.emitInstruction(instruction.OpSub),
	)
}
func (c *Compiler) scopeSM(node ast.Expression) func() {
	c.scope().code.sm.Push(c.pos(), node.Location(), reflect.TypeOf(node).String())
	return func() {
		c.scope().code.sm.Pop(c.pos())
	}
}
func (c *Compiler) compileWhileExpression(node ast.WhileExpression) error {
	whileCode, err := c.makeClosure(func() (*code, error) { // should this really be a closure?
		return c.makecb(func() error {
			defer c.scopeSM(node)()
			return c.compileWhileExpressionBody(node)
		})
	}, nil, object.CodeReturnScopeLoop)
	return iferr(
		err,
		c.annotate("outer while start"),
		c.emit(whileCode),
		c.emitInstruction(instruction.OpCall, 0),
	)
}
func (c *Compiler) compileForExpression(node ast.ForExpression) error {
	// convert the for into a while:
	/**

	regular loop:
	let __ii = iteritems(arr);
	let __l = len(ii);
	let __i = -1;
	let __res = null; // resolve to this
	while __i++ < __l {
		let key = __ii.(__i).0;
		let value = __ii.(__i).1;

		// body
	};

	loop with an arrow expression:
	let __ii = iteritems(arr);
	let __l = len(ii);
	let __i = -1;
	let __r = []; // resolve to this
	while __i++ < __l {
		let key = __ii.(__i).0;
		let value = __ii.(__i).1;

		__r = __r + [<body()>];
	};

	*/

	_, isArrow := node.Body.(ast.ArrowExpression)

	var whileBody []ast.Statement
	if node.Index != nil {
		whileBody = append(whileBody, ast.Statement{Expr: ast.LetExpression{
			Identifiers: []ast.Identifier{*node.Index},
			Initialization: ast.FieldAccessExpression{
				Left: ast.FieldAccessExpression{
					Left:  ast.Identifier{Name: "!ii"},
					Right: ast.Identifier{Name: "!i"},
				},
				Right: ast.NumberExpression{Value: 0},
			},
		}})
	}
	whileBody = append(whileBody, ast.Statement{Expr: ast.LetExpression{
		Identifiers: []ast.Identifier{node.Value},
		Initialization: ast.FieldAccessExpression{
			Left: ast.FieldAccessExpression{
				Left:  ast.Identifier{Name: "!ii"},
				Right: ast.Identifier{Name: "!i"},
			},
			Right: ast.NumberExpression{Value: 1},
		},
	}})
	if isArrow {
		//whileBody = append(whileBody, ast.Statement{Expr: ast.AssignExpression{
		//	Identifier: ast.Identifier{Name: "!r"},
		//	Value: ast.PlusExpression{
		//		Left:  ast.Identifier{Name: "!r"},
		//		Right: ast.ArrayExpression{Items: []ast.Expression{node.Body.(ast.ArrowExpression).Expr}},
		//	},
		//}})
		whileBody = append(whileBody, ast.Statement{Expr: ast.CallExpression{
			Callee:    ast.Identifier{Name: "append"},
			Arguments: []ast.Expression{ast.Identifier{Name: "!r"}, node.Body.(ast.ArrowExpression).Expr},
		}})
	} else {
		whileBody = append(whileBody, ast.Statement{Expr: node.Body})
	}

	while := ast.BlockExpression{
		Stmts: []ast.Statement{
			{Expr: ast.LetExpression{ // let ii = iteritems(arr);
				Identifiers: []ast.Identifier{{Name: "!ii"}}, // ! is added to the var name to guarantee that it does not collide with user-specified local variables
				Initialization: ast.CallExpression{
					Callee:    ast.Identifier{Name: "iteritems"},
					Arguments: []ast.Expression{node.Range},
				},
			}},
			{Expr: ast.LetExpression{ // let l = len(ii);
				Identifiers: []ast.Identifier{{Name: "!l"}},
				Initialization: ast.CallExpression{
					Callee:    ast.Identifier{Name: "len"},
					Arguments: []ast.Expression{ast.Identifier{Name: "!ii"}},
				},
			}},
			{Expr: ast.LetExpression{ // let i = -1;
				Identifiers:    []ast.Identifier{{Name: "!i"}},
				Initialization: ast.NumberExpression{Value: -1},
			}},
			{Expr: ast.LetExpression{ // let r = [];
				Identifiers:    []ast.Identifier{{Name: "!r"}},
				Initialization: ast.ArrayExpression{Items: nil},
			}},
			{Expr: ast.WhileExpression{ // while __i++ < __l {
				Loc: node.Location(),
				Condition: ast.LtExpression{
					Left: ast.AssignExpression{
						Identifier: ast.Identifier{Name: "!i"},
						Value: ast.PlusExpression{
							Left:  ast.Identifier{Name: "!i"},
							Right: ast.NumberExpression{Value: 1},
						},
					},
					Right: ast.Identifier{Name: "!l"},
				},
				Body: ast.BlockExpression{
					Stmts: whileBody,
				},
			}},
		},
	}

	c.pushSymbolsLinked()
	whileCode, err := c.make(while)
	sym := c.popSymbols()
	if err != nil {
		return err
	}

	c.annotate("for start")
	err = c.emit(whileCode)
	if isArrow {
		err = iferr(err,
			c.emitInstruction(instruction.OpPop), // remove the "null"
			c.emitPushSymbol(sym.getLocal("!r")),
		)
	}
	return err
}
func (c *Compiler) compileLetExpression(node ast.LetExpression) error {
	syms := make([]*Symbol, len(node.Identifiers))
	for i, id := range node.Identifiers {
		if c.symbols.hasLocal(id.Name) {
			return fmt.Errorf("identifier already declared in this scope: %s", id.Name)
		}

		syms[i] = c.symbols.createLocal(id.Name)
	}

	var err error
	err = iferr(
		c.emitNode(node.Initialization),
	)

	if len(syms) == 1 {
		err = iferr(
			err,
			c.emitStoreSymbol(syms[0]),
		)
	} else {
		err = iferr(
			err,
			c.emitInstruction(instruction.OpDup),
			c.emitInstruction(instruction.OpUntuple, len(syms)),
		)
		for _, sym := range syms {
			err = iferr(
				err,
				c.emitStoreSymbol(sym),
				c.emitInstruction(instruction.OpPop),
			)
		}
	}

	return err
}
func (c *Compiler) compileTupleAssignExpression(node ast.TupleAssignExpression) error {
	var err error
	err = iferr(
		c.emitNode(node.Value),
		c.emitInstruction(instruction.OpDup),
		c.emitInstruction(instruction.OpUntuple, len(node.Tuple.Exprs)),
	)

	for _, expr := range node.Tuple.Exprs {
		switch expr := expr.(type) {
		case ast.Identifier:
			sym, _ := c.symbols.get(expr.Name)
			if sym == nil {
				return fmt.Errorf("identifier is not declared in this scope: " + expr.Name)
			}
			err = iferr(err,
				c.emitStoreSymbol(sym),
				c.emitInstruction(instruction.OpPop),
			)
		case ast.FieldAccessExpression:
			err = iferr(err,
				c.emitNode(expr.Right),
				c.emitNode(expr.Left),
				c.emitInstruction(instruction.OpFieldAssign),
				c.emitInstruction(instruction.OpPop),
			)
		default:
			return fmt.Errorf("identifier or dot-access expected on the lval tuple, got: " + reflect.TypeOf(expr).String())
		}
	}

	return err
}
func (c *Compiler) compileAssignExpression(node ast.AssignExpression) error {
	sym, _ := c.symbols.get(node.Identifier.Name) // todo: this can resolve to outer scope
	if sym == nil {
		return fmt.Errorf("cannot assign to unknown identifier: %s", node.Identifier.Name)
	}

	return iferr(
		c.emitNode(node.Value),
		c.emitStoreSymbol(sym),
	)
}
func (c *Compiler) compileIdentifier(node ast.Identifier) error {
	s, _ := c.symbols.get(node.Name)
	if s == nil {
		return fmt.Errorf("unknown identifier: %s", node.Name)
	}

	return c.emitPushSymbol(s)
}
func (c *Compiler) compileCallExpression(node ast.CallExpression) error {
	return iferr(
		c.compileCopiedExpressions(node.Arguments),
		c.emitNode(node.Callee),
		c.emitInstruction(instruction.OpCall, len(node.Arguments)),
	)
}
func (c *Compiler) compileReturnExpression(node ast.ReturnExpression) error {
	// todo: check if inside a func
	return iferr(
		c.emitNode(node.Expr),
		c.emitInstruction(instruction.OpReturn, int(object.CodeReturnScopeFunc)),
	)
}
func (c *Compiler) compileBreakExpression(node ast.BreakExpression) error {
	// todo: check if inside a loop
	return iferr(
		c.emitPushNull(), // loop resolves to null when break'ed
		c.emitInstruction(instruction.OpReturn, int(object.CodeReturnScopeLoop)),
	)
}
func (c *Compiler) compileContinueExpression(node ast.ContinueExpression) error {
	// todo: check if inside a loop?
	return iferr(
		c.emitPushNull(),
		c.emitInstruction(instruction.OpJmp, int(LabelAddress), int(instruction.LabelKindContinue)),
	)
}
func (c *Compiler) compileArrowExpression(node ast.ArrowExpression) error {
	// todo: can be different depending on the context
	return iferr(
		c.emitNode(node.Expr),
	)
}
func (c *Compiler) compileGroupExpression(node ast.GroupExpression) error {
	return c.emitNode(node.Expr)
}
func (c *Compiler) compileStatement(node ast.Statement) error {
	return iferr(
		c.emitNode(node.Expr),
		c.emitInstruction(instruction.OpPop),
	)
}
func (c *Compiler) compileBlockExpression(node ast.BlockExpression) error {
	// todo: scope change? frame?
	return iferr(
		c.compileStatements(node.Stmts),
		c.emitPushNull(), // todo: block expressions always resolve to null (at least for now). Maybe allow to break from block expressions?
	)
}
func (c *Compiler) compileExports(node ast.Exports) error {
	var err error
	for field, expr := range node.Fields {
		if expr != nil {
			err = iferr(err, c.compileLetExpression(ast.LetExpression{
				Identifiers:    []ast.Identifier{{Name: field}},
				Initialization: expr,
			}), c.emitInstruction(instruction.OpPop))
		}
	}

	return iferr(
		err,
		c.emitPushNull(),
	)
}
func (c *Compiler) compileImport(node ast.Import) error {
	return iferr(
		c.emitNode(node.Module),
		c.emitInstruction(instruction.OpImport), // resolves to code
		c.emitInstruction(instruction.OpClosure),
		c.emitInstruction(instruction.OpCall, 0),
	)
}
func (c *Compiler) CompileModule(node ast.Module) (*Module, error) {
	if len(c.scopes) > 1 || c.symbols.parent != nil {
		return nil, fmt.Errorf("unexpected module compilation on non-root level")
	}

	c.enterScope()
	c.pushSymbols()
	c.annotate("module: " + node.Name)

	var exports *ast.Exports
	var err error
	for _, stmt := range node.Block.Stmts {
		if _, ok := stmt.Expr.(ast.Exports); ok {
			tmp := stmt.Expr.(ast.Exports)
			exports = &tmp
		}

		err = iferr(err, c.emitNode(stmt))
	}

	if exports != nil && len(exports.Fields) > 0 {
		for field := range exports.Fields {
			symbol := c.symbols.getLocal(field)
			if symbol == nil {
				return nil, fmt.Errorf("exported field declaration is missing: %s", field)
			}

			err = iferr(
				err,
				c.emitPushSymbol(symbol),
				c.emitConstantObject(&object.String{Value: field}),
			)
		}
		err = iferr(err, c.emitInstruction(instruction.OpStruct, len(exports.Fields)))
	} else {
		err = iferr(err, c.emitInstruction(instruction.OpStruct, 0))
	}

	s := c.leaveScope()
	sym := c.popSymbols()
	if err != nil {
		return nil, err
	}

	if len(c.symbols.foreign) > 0 {
		return nil, fmt.Errorf("module is not expected to have foreigns")
	}

	id := c.registerObject(&object.Code{
		Code:        s.code.b,
		Locals:      *sym.locals,
		ReturnScope: object.CodeReturnScopeFunc,
	})
	c.saveDebugData(id, s.code.sm, sym)

	module := &Module{
		EntryPoint: uint16(id),
		Objects:    c.objects,
		DebugData:  c.debugData,
	}

	return module, nil
}
func (c *Compiler) compileFuncExpression(node ast.FuncExpression) error {
	fc, err := c.makeClosure(func() (*code, error) {
		return c.makecb(func() error {
			return iferr(
				//c.annotate("func at "+node.Location().String()),
				c.emitNode(node.Body),
			)
		})
	}, node.Arguments, object.CodeReturnScopeFunc)
	return iferr(
		err,
		c.emit(fc),
	)
}

func (c *Compiler) emitNode(node ast.Expression) error {
	if node == nil {
		return nil
	}

	defer c.scopeSM(node)()

	//c.annotate(node.Location().String() + ": " + node.String())
	switch node := node.(type) {
	case ast.NumberExpression:
		return c.compileNumberExpression(node)
	case ast.ArrayExpression:
		return c.compileArrayExpression(node)
	case ast.TupleExpression:
		return c.compileTupleExpression(node)
	case ast.FieldAccessExpression:
		return c.compileFieldAccessExpression(node)
	case ast.FieldAssignExpression:
		return c.compileFieldAssignExpression(node)
	case ast.String:
		return c.compileString(node)
	case ast.PlusExpression:
		return c.compilePlusExpression(node)
	case ast.MinusExpression:
		return c.compileMinusExpression(node)
	case ast.MultExpression:
		return c.compileMultExpression(node)
	case ast.DivExpression:
		return c.compileDivExpression(node)
	case ast.ModExpression:
		return c.compileModExpression(node)
	case ast.GtExpression:
		return c.compileGtExpression(node)
	case ast.LtExpression:
		return c.compileLtExpression(node)
	case ast.GteExpression:
		return c.compileGteExpression(node)
	case ast.LteExpression:
		return c.compileLteExpression(node)
	case ast.LogicalAndExpression:
		return c.compileLogicalAndExpression(node)
	case ast.LogicalOrExpression:
		return c.compileLogicalOrExpression(node)
	case ast.EqTestExpression:
		return c.compileEqTestExpression(node)
	case ast.Module:
		panic("unexpected module compilation")
	case ast.IfExpression:
		return c.compileIfExpression(node)
	case ast.LetExpression:
		return c.compileLetExpression(node)
	case ast.Identifier:
		return c.compileIdentifier(node)
	case ast.BlockExpression:
		return c.compileBlockExpression(node)
	case ast.FuncExpression:
		return c.compileFuncExpression(node)
	case ast.AssignExpression:
		return c.compileAssignExpression(node)
	case ast.CallExpression:
		return c.compileCallExpression(node)
	case ast.ReturnExpression:
		return c.compileReturnExpression(node)
	case ast.BreakExpression:
		return c.compileBreakExpression(node)
	case ast.ContinueExpression:
		return c.compileContinueExpression(node)
	case ast.ArrowExpression:
		return c.compileArrowExpression(node)
	case ast.GroupExpression:
		return c.compileGroupExpression(node)
	case ast.Statement:
		return c.compileStatement(node)
	case ast.WhileExpression:
		return c.compileWhileExpression(node)
	case ast.Exports:
		return c.compileExports(node)
	case ast.Import:
		return c.compileImport(node)
	case ast.NegationExpression:
		return c.compileNegationExpression(node)
	case ast.PrefixMinusExpression:
		return c.compilePrefixMinusExpression(node)
	case ast.ForExpression:
		return c.compileForExpression(node)
	case ast.StructExpression:
		return c.compileStructExpression(node)
	case ast.MapExpression:
		return c.compileMapExpression(node)
	case ast.TupleAssignExpression:
		return c.compileTupleAssignExpression(node)

	default:
		panic("dont know how to compile this type: " + reflect.TypeOf(node).String())
	}
}
