package eval

import (
	"fmt"
	"os"
	"reflect"
	"ryanlang/ast"
	"ryanlang/funcs"
	"ryanlang/lexer"
	"ryanlang/object"
	"ryanlang/parser"
	"strconv"
)

func (e *Evaluator) evalPlusExpression(expr ast.PlusExpression) object.Object {
	return funcs.Plus(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalMinusExpression(expr ast.MinusExpression) object.Object {
	return funcs.Minus(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalModExpression(expr ast.ModExpression) object.Object {
	return funcs.Mod(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalPrefixMinusExpression(expr ast.PrefixMinusExpression) object.Object {
	var right object.Object

	if right = e.expectEvalToAnyType(expr.Expr); object.IsError(right) {
		return right
	}

	if right.Type() == object.NUMBER {
		return &object.Number{Value: -right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("don't know how to negate type: %s", right.Type().String())}
	}
}
func (e *Evaluator) evalMultExpression(expr ast.MultExpression) object.Object {
	return funcs.Mult(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalDivExpression(expr ast.DivExpression) object.Object {
	return funcs.Div(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalLogicalAndExpression(expr ast.LogicalAndExpression) object.Object {
	// key point is to first eval left, and only eval right is left is false, do not eval right otherwise

	var left, right object.Object

	if left = e.expectEvalToType(expr.Left, object.BOOLEAN); object.IsError(left) {
		return left
	}
	if !left.(*object.Boolean).Value { // left is false, no need to evaluate right
		return &object.Boolean{Value: false}
	}

	if right = e.expectEvalToType(expr.Right, object.BOOLEAN); object.IsError(right) {
		return right
	}

	return &object.Boolean{Value: right.(*object.Boolean).Value}
}
func (e *Evaluator) evalLogicalOrExpression(expr ast.LogicalOrExpression) object.Object {
	return funcs.LogicalOr(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalLetExpression(expr ast.LetExpression) object.Object {
	var inits []object.Object
	var ret object.Object

	if len(expr.Identifiers) > 1 {
		if init := e.expectEvalToType(expr.Initialization, object.TUPLE); object.IsError(init) {
			return init
		} else {
			inits = init.(*object.Tuple).Values
			ret = init
		}
	} else {
		if init := e.expectEvalToAnyType(expr.Initialization); object.IsError(init) {
			return init
		} else {
			inits = append(inits, init)
			ret = init
		}
	}

	if len(expr.Identifiers) != len(inits) {
		return &object.Error{Msg: "values expected on right side: " + strconv.Itoa(len(expr.Identifiers))}
	}

	for index, id := range expr.Identifiers {
		_, ok := e.env.GetFromCurrent(id.Name)
		if ok {
			return &object.Error{Msg: "identifier already declared: " + id.Name, Loc: id.Location()}
		}

		e.env.Set(id.Name, inits[index])
	}

	return ret
}
func (e *Evaluator) evalIdentifier(expr ast.Identifier) object.Object {
	value, ok := e.env.Get(expr.Name)
	if !ok {
		return &object.Error{Msg: "undefined identifier: " + expr.Name, Loc: expr.Location()}
	}
	return value
}
func (e *Evaluator) evalString(expr ast.String) object.Object {
	return &object.String{Value: expr.Value}
}
func (e *Evaluator) evalNumber(expr ast.NumberExpression) object.Object {
	return &object.Number{Value: expr.Value}
}
func (e *Evaluator) assignExpression(id ast.Identifier, value object.Object) object.Object {
	currentValue, ok := e.env.Get(id.Name)
	if !ok {
		return &object.Error{
			Msg: "undeclared identifier",
		}
	}
	// todo: check if it's a builtin symbol, and probably disallow assigning to builtin symbols?
	if value.Type() != currentValue.Type() {
		return &object.Error{Msg: "identifier already holds value of type: " + currentValue.Type().String(), Loc: id.Loc}
	}
	e.env.Replace(id.Name, value)
	return value
}
func (e *Evaluator) evalAssignExpression(expr ast.AssignExpression) object.Object {
	var value object.Object
	if value = e.expectEvalToAnyType(expr.Value); object.IsError(value) {
		return value
	}
	return e.assignExpression(expr.Identifier, value)
}
func (e *Evaluator) fieldAssignExpressionValue(fa ast.FieldAccessExpression, value object.Object) object.Object {
	return funcs.FieldAssign(e.expectEvalToAnyType(fa.Left), e.expectEvalToAnyType(fa.Right), value)
}
func (e *Evaluator) evalFieldAssignExpression(expr ast.FieldAssignExpression) object.Object {
	var value object.Object
	if value = e.expectEvalToAnyType(expr.Value); object.IsError(value) {
		return value
	}

	return e.fieldAssignExpressionValue(expr.FieldAccess, value)
}
func (e *Evaluator) evalFuncExpression(expr ast.FuncExpression) object.Object {
	return &object.Function{Node: expr, Env: e.env}
}
func (e *Evaluator) evalBlockExpression(expr ast.BlockExpression) object.Object {
	var ret object.Object
	derivedEvaluator := NewWithEnv(e.env.Derive())
	for _, stmt := range expr.Stmts {
		if ret = derivedEvaluator.expectEvalToAnyType(stmt); object.IsError(ret) {
			return ret
		}
		if ret.Type() == object.RETURNOBJECT {
			return ret
		}
		if ret.Type() == object.BREAK {
			return ret
		}
		if ret.Type() == object.CONTINUE {
			return ret
		}
	}
	return &object.Null{} // todo: maybe allow block expression to resolve to something?
}
func (e *Evaluator) evalModule(expr ast.Module) object.Object {
	var exports *object.Exports
	derivedEvaluator := NewWithEnv(e.env.Derive())
	for _, stmt := range expr.Block.Stmts {
		var ret object.Object
		if ret = derivedEvaluator.expectEvalToAnyType(stmt.Expr); object.IsError(ret) {
			return ret
		}
		if ret.Type() == object.RETURNOBJECT {
			return &object.Error{Msg: "cannot return from a top-level in a module"}
		}
		if ret.Type() == object.BREAK {
			return &object.Error{Msg: "cannot break from a top-level in a module"}
		}
		if ret.Type() == object.CONTINUE {
			return &object.Error{Msg: "cannot continue in a top-level in a module"}
		}
		if ret.Type() == object.EXPORTS {
			exports = ret.(*object.Exports)
		}
	}
	ret := &object.Module{
		Name:    expr.Name,
		Exports: map[string]object.Object{},
	}
	if exports != nil {
		for _, field := range exports.Fields {
			v, ok := derivedEvaluator.env.GetFromCurrent(field)
			if !ok {
				return &object.Error{Msg: "exported identifier is not declared: " + field}
			}
			ret.Exports[field] = v
		}
	}
	return ret
}
func (e *Evaluator) evalImport(expr ast.Import) object.Object {
	var fn object.Object

	if fn = e.expectEvalToType(expr.Module, object.STRING); object.IsError(fn) {
		return fn
	}

	f, err := os.Open(fn.(*object.String).Value)
	if err != nil {
		return &object.Error{Msg: err.Error()}
	}
	l := lexer.New(f, fn.(*object.String).Value)
	p := parser.New(l)

	mod := p.ReadModule(fn.(*object.String).Value)
	return New().Eval(mod)
}
func (e *Evaluator) evalGroupExpression(expr ast.GroupExpression) object.Object {
	return e.expectEvalToAnyType(expr.Expr)
}
func (e *Evaluator) evalStatement(expr ast.Statement) object.Object {
	return e.expectEvalToAnyType(expr.Expr)
}
func (e *Evaluator) evalCallExpression(expr ast.CallExpression) object.Object {
	var callee object.Object
	if callee = e.expectEvalToType(expr.Callee, object.FUNCTION); object.IsError(callee) {
		return callee
	}

	funcArguments := callee.(*object.Function).Node.Arguments
	if len(funcArguments) != len(expr.Arguments) {
		return &object.Error{Msg: fmt.Sprintf("expected %d arguments, got %d", len(funcArguments), len(expr.Arguments)), Loc: expr.Callee.Location()}
	}

	derivedEvaluator := NewWithEnv(callee.(*object.Function).Env.Derive())
	for i, argExpr := range expr.Arguments {
		var argValue object.Object
		if argValue = e.expectEvalToAnyType(argExpr); object.IsError(argValue) {
			return argValue
		}

		derivedEvaluator.env.Set(funcArguments[i].Name, argValue)
	}

	var ret object.Object
	if ret = derivedEvaluator.expectEvalToAnyType(callee.(*object.Function).Node.Body); object.IsError(ret) {
		return ret
	}

	// todo: disallow break/continue outside of their scopes
	// i.e. now it's possible to do a break or continue inside a function
	// even though there's no enclosing loop cycle
	// e.g. func { break; }
	if ret.Type() != object.RETURNOBJECT {
		//return &object.Error{Msg: "missing return"}
		return &object.Null{}
	}

	return ret.(*object.ReturnObject).Obj
}
func (e *Evaluator) evalReturnExpression(expr ast.ReturnExpression) object.Object {
	var ret object.Object
	if ret = e.expectEvalToAnyType(expr.Expr); object.IsError(ret) {
		return ret
	}
	return &object.ReturnObject{Obj: ret}
}
func (e *Evaluator) evalContinueExpression(expr ast.ContinueExpression) object.Object {
	return &object.ContinueObject{}
}
func (e *Evaluator) evalBreakExpression(expr ast.BreakExpression) object.Object {
	return &object.BreakObject{}
}
func (e *Evaluator) evalIfExpression(expr ast.IfExpression) object.Object {
	var condition object.Object

	derivedEvaluator := NewWithEnv(e.env.Derive())
	// condition should be eval'ed inside a derived env, because:
	/**
	let f = func=>10;
	if (let x = f()) > 5 {
		println(x);
	};
	println(x); // should be undefined
	*/
	if condition = derivedEvaluator.expectEvalToType(expr.Condition, object.BOOLEAN); object.IsError(condition) {
		return condition
	}

	var ret object.Object
	if condition.(*object.Boolean).Value {
		if ret = derivedEvaluator.expectEvalToAnyType(expr.Then); object.IsError(ret) {
			return ret
		}
	} else if expr.Else != nil {
		if ret = derivedEvaluator.expectEvalToAnyType(expr.Else); object.IsError(ret) {
			return ret
		}
	} else {
		ret = &object.StaticNull
	}

	return ret
}
func (e *Evaluator) evalGtExpression(expr ast.GtExpression) object.Object {
	return funcs.Gt(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalLtExpression(expr ast.LtExpression) object.Object {
	var left, right object.Object

	if left = e.expectEvalToType(expr.Left, object.NUMBER); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToType(expr.Right, object.NUMBER); object.IsError(right) {
		return right
	}

	return &object.Boolean{Value: left.(*object.Number).Value < right.(*object.Number).Value}
}
func (e *Evaluator) evalGteExpression(expr ast.GteExpression) object.Object {
	var left, right object.Object

	if left = e.expectEvalToType(expr.Left, object.NUMBER); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToType(expr.Right, object.NUMBER); object.IsError(right) {
		return right
	}

	return &object.Boolean{Value: left.(*object.Number).Value >= right.(*object.Number).Value}
}
func (e *Evaluator) evalLteExpression(expr ast.LteExpression) object.Object {
	var left, right object.Object

	if left = e.expectEvalToType(expr.Left, object.NUMBER); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToType(expr.Right, object.NUMBER); object.IsError(right) {
		return right
	}

	return &object.Boolean{Value: left.(*object.Number).Value <= right.(*object.Number).Value}
}
func (e *Evaluator) evalBuiltinFunction(expr ast.BuiltinFunction) object.Object {
	if _, ok := funcs.BuiltinFunctions[expr.Name]; !ok {
		panic("unknown built-in function: " + expr.Name)
	}
	args := map[string]object.Object{}
	for _, argName := range funcs.BuiltinFunctions[expr.Name].Arguments {
		v, ok := e.env.Get(argName)
		if !ok {
			return &object.Error{Msg: fmt.Sprintf("`%s` parameter required", argName)}
		}
		args[argName] = v
	}
	return funcs.BuiltinFunctions[expr.Name].Body(args)
}
func (e *Evaluator) evalStructExpression(expr ast.StructExpression) object.Object {
	ret := &object.Struct{
		Fields: map[string]object.Object{},
	}

	//newEnv := object.NewEnvironment()

	for field, fieldExpr := range expr.Fields {
		var fieldVal object.Object
		if fieldVal = e.expectEvalToAnyType(fieldExpr); object.IsError(fieldVal) {
			return fieldVal
		}
		ret.Fields[field] = fieldVal
		//newEnv.Set(field, fieldVal)
	}

	// add "this" to all functions in the struct
	for _, fieldVal := range ret.Fields {
		if fieldVal.Type() == object.FUNCTION {
			fieldVal.(*object.Function).Env = fieldVal.(*object.Function).Env.DeriveWith("this", ret)
		}
	}

	return ret
}
func (e *Evaluator) evalMapExpression(expr ast.MapExpression) object.Object {
	ret := &object.Map{
		Fields: map[string]object.MapItem{},
	}

	//newEnv := object.NewEnvironment()

	for _, mapField := range expr.Fields {
		var field object.Object
		if field = e.expectEvalToHashableType(mapField.Key); object.IsError(field) {
			return field
		}

		var fieldVal object.Object
		if fieldVal = e.expectEvalToAnyType(mapField.Value); object.IsError(fieldVal) {
			return fieldVal
		}
		ret.Fields[field.(object.Hashable).Hash()] = object.MapItem{
			Key:   field,
			Value: fieldVal,
		}
		//newEnv.Set(field, fieldVal)
	}

	return ret
}
func (e *Evaluator) evalExports(expr ast.Exports) object.Object {
	ret := &object.Exports{}

	//newEnv := object.NewEnvironment()

	for field, fieldExpr := range expr.Fields {
		if fieldExpr != nil { // "nil" means this field is declared later in the module body
			var fieldVal object.Object
			if fieldVal = e.expectEvalToAnyType(fieldExpr); object.IsError(fieldVal) {
				return fieldVal
			}
			e.env.Set(field, fieldVal)
		}
		ret.Fields = append(ret.Fields, field)

		//newEnv.Set(field, fieldVal)
	}

	return ret
}
func (e *Evaluator) evalFieldAccessExpression(expr ast.FieldAccessExpression) object.Object {
	return funcs.FieldAccess(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalTupleAssignExpression(expr ast.TupleAssignExpression) object.Object {
	var rval object.Object
	if rval = e.expectEvalToType(expr.Value, object.TUPLE); object.IsError(rval) {
		return rval
	}

	if len(expr.Tuple.Exprs) != len(rval.(*object.Tuple).Values) {
		return &object.Error{Msg: "rval tuple has wrong number of values: " + strconv.Itoa(len(rval.(*object.Tuple).Values)) + ", expected: " + strconv.Itoa(len(expr.Tuple.Exprs))}
	}

	var lval object.Object
	if lval = e.expectEvalToType(expr.Tuple, object.TUPLE); object.IsError(lval) {
		return lval
	}
	for index, exp := range expr.Tuple.Exprs {
		switch exp.(type) {
		case ast.Identifier:
			e.assignExpression(exp.(ast.Identifier), rval.(*object.Tuple).Values[index])
		case ast.FieldAccessExpression:
			e.fieldAssignExpressionValue(exp.(ast.FieldAccessExpression), rval.(*object.Tuple).Values[index])
		default:
			return &object.Error{Msg: "identifier or dot-access expected on the lval tuple, got: " + reflect.TypeOf(exp).String()}
		}
	}

	return &object.Null{}
}
func (e *Evaluator) evalWhileExpression(expr ast.WhileExpression) object.Object {
	var condition object.Object
	var ret object.Object = &object.StaticNull

	for {
		if condition = e.expectEvalToType(expr.Condition, object.BOOLEAN); object.IsError(condition) {
			return condition
		}

		if !condition.(*object.Boolean).Value {
			break
		}

		if ret = e.expectEvalToAnyType(expr.Body); object.IsError(ret) {
			return ret
		}
		if ret.Type() == object.RETURNOBJECT {
			break
		}
		if ret.Type() == object.BREAK {
			ret = &object.Null{}
			break
		}
		if ret.Type() == object.CONTINUE {
			ret = &object.Null{}
			// nothing to do, just continue
		}

		// todo: support arrow expression?
	}

	return ret
}
func (e *Evaluator) evalForExpression(expr ast.ForExpression) object.Object {
	var r object.Object
	var ret object.Object = &object.StaticNull

	if r = e.expectEvalToAnyType(expr.Range); object.IsError(r) {
		return r
	}

	var rangeItems []struct {
		index object.Object
		value object.Object
	}
	if r.Type() == object.ARRAY {
		for k, v := range r.(*object.Array).Items {
			rangeItems = append(rangeItems, struct {
				index object.Object
				value object.Object
			}{index: &object.Number{Value: k}, value: v})
		}
	} else if r.Type() == object.MAP {
		for _, item := range r.(*object.Map).Fields {
			rangeItems = append(rangeItems, struct {
				index object.Object
				value object.Object
			}{index: item.Key, value: item.Value})
		}
	} else {
		return &object.Error{Msg: "cannot iterate over type " + r.Type().String()}
	}

	_, returnArrowExpression := expr.Body.(ast.ArrowExpression)
	arrowItems := &object.Array{}

	derivedEvaluator := NewWithEnv(e.env.Derive())
	for _, v := range rangeItems {
		if expr.Index != nil {
			derivedEvaluator.env.Set(expr.Index.Name, v.index)
		}
		derivedEvaluator.env.Set(expr.Value.Name, v.value)

		if ret = derivedEvaluator.expectEvalToAnyType(expr.Body); object.IsError(ret) {
			return ret
		}
		if ret.Type() == object.RETURNOBJECT {
			returnArrowExpression = false
			break
		}
		if ret.Type() == object.BREAK {
			ret = &object.Null{}
			break
		}
		if ret.Type() == object.CONTINUE {
			ret = &object.Null{}
			// nothing to do, just continue
			continue
		}
		if returnArrowExpression {
			arrowItems.Items = append(arrowItems.Items, ret)
			ret = &object.Null{}
		}
	}
	if returnArrowExpression {
		ret = arrowItems
	}

	return ret
}
func (e *Evaluator) evalArrayExpression(expr ast.ArrayExpression) object.Object {
	items := make([]object.Object, len(expr.Items))
	for i, item := range expr.Items {
		var val object.Object
		if val = e.expectEvalToAnyType(item); object.IsError(val) {
			return val
		}

		items[i] = val
	}
	return &object.Array{Items: items}
}
func (e *Evaluator) evalNegationExpression(expr ast.NegationExpression) object.Object {
	var val object.Object
	if val = e.expectEvalToType(expr.Expr, object.BOOLEAN); object.IsError(val) {
		return val
	}
	return object.StaticBool(!val.(*object.Boolean).Value)
}
func (e *Evaluator) evalEqTestExpression(expr ast.EqTestExpression) object.Object {
	return funcs.EqTest(e.expectEvalToAnyType(expr.Left), e.expectEvalToAnyType(expr.Right))
}
func (e *Evaluator) evalArrowExpression(expr ast.ArrowExpression) object.Object {
	return e.expectEvalToAnyType(expr.Expr)
}
func (e *Evaluator) evalTupleExpression(expr ast.TupleExpression) object.Object {
	ret := &object.Tuple{}

	for _, exp := range expr.Exprs {
		var v object.Object
		if v = e.expectEvalToAnyType(exp); object.IsError(v) {
			return v
		}
		ret.Values = append(ret.Values, v)
	}
	return ret
}
