package eval

import (
	"fmt"
	"os"
	"reflect"
	"ryanlang/ast"
	"ryanlang/lexer"
	"ryanlang/object"
	"ryanlang/parser"
	"strconv"
)

func (e *Evaluator) evalPlusExpression(expr ast.PlusExpression) object.Object {
	var left, right object.Object

	if left = e.expectEvalToAnyType(expr.Left); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToAnyType(expr.Right); object.IsError(right) {
		return right
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value + right.(*object.Number).Value}
	} else if left.Type() == object.STRING && right.Type() == object.STRING {
		return &object.String{Value: left.(*object.String).Value + right.(*object.String).Value}
	} else if left.Type() == object.ARRAY && right.Type() == object.ARRAY {
		ret := append([]object.Object{}, left.(*object.Array).Items...)
		return &object.Array{Items: append(ret, right.(*object.Array).Items...)}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for a plus operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func (e *Evaluator) evalMinusExpression(expr ast.MinusExpression) object.Object {
	var left, right object.Object

	if left = e.expectEvalToAnyType(expr.Left); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToAnyType(expr.Right); object.IsError(right) {
		return right
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value - right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for a minus operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func (e *Evaluator) evalModExpression(expr ast.ModExpression) object.Object {
	var left, right object.Object

	if left = e.expectEvalToAnyType(expr.Left); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToAnyType(expr.Right); object.IsError(right) {
		return right
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value % right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for a mod operator: %s, %s", left.Type().String(), right.Type().String())}
	}
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
	var left, right object.Object

	if left = e.expectEvalToAnyType(expr.Left); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToAnyType(expr.Right); object.IsError(right) {
		return right
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value * right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for a mult operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func (e *Evaluator) evalDivExpression(expr ast.DivExpression) object.Object {
	var left, right object.Object

	if left = e.expectEvalToAnyType(expr.Left); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToAnyType(expr.Right); object.IsError(right) {
		return right
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return &object.Number{Value: left.(*object.Number).Value / right.(*object.Number).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for a div operator: %s, %s", left.Type().String(), right.Type().String())}
	}
}
func (e *Evaluator) evalLogicalAndExpression(expr ast.LogicalAndExpression) object.Object {
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
	var left, right object.Object

	if left = e.expectEvalToAnyType(expr.Left); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToAnyType(expr.Right); object.IsError(right) {
		return right
	}

	if left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN {
		return &object.Boolean{Value: left.(*object.Boolean).Value || right.(*object.Boolean).Value}
	} else {
		return &object.Error{Msg: fmt.Sprintf("incompatible types for a logical or operator: %s, %s", left.Type().String(), right.Type().String())}
	}
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
func (e *Evaluator) evalNull(expr ast.Null) object.Object {
	return &object.Null{}
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
	var lval object.Object

	if lval = e.expectEvalToAnyType(fa.Left); object.IsError(lval) {
		return lval
	}

	var rval object.Object
	if lval.Type() == object.STRUCT {
		if rval = e.expectEvalToType(fa.Right, object.STRING); object.IsError(rval) {
			return rval
		}
		fieldName := rval.(*object.String).Value
		currentValue, ok := lval.(*object.Struct).Fields[fieldName]
		if !ok {
			return &object.Error{
				Msg: "cannot assign to a non-existing field: " + fieldName,
				Loc: fa.Right.Location(),
			}
		}
		if currentValue.Type() != value.Type() {
			return &object.Error{
				Msg: "field already holds a value of type " + currentValue.Type().String() + ", got: " + value.Type().String(),
			}
		}
		lval.(*object.Struct).Fields[fieldName] = value
	} else if lval.Type() == object.MAP {
		if rval = e.expectEvalToHashableType(fa.Right); object.IsError(rval) {
			return rval
		}
		hash := rval.(object.Hashable).Hash()
		currentValue, ok := lval.(*object.Map).Fields[hash]
		if ok && currentValue.Type() != value.Type() {
			return &object.Error{
				Msg: "field already holds a value of type " + currentValue.Type().String() + ", got: " + value.Type().String(),
			}
		}
		lval.(*object.Map).Fields[hash] = value
	} else if lval.Type() == object.ARRAY {
		if rval = e.expectEvalToType(fa.Right, object.NUMBER); object.IsError(rval) {
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
	for _, expr := range expr.Exprs {
		if ret = derivedEvaluator.expectEvalToAnyType(expr); object.IsError(ret) {
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
	return &object.Null{}
}
func (e *Evaluator) evalModule(expr ast.Module) object.Object {
	var exports *object.Exports
	derivedEvaluator := NewWithEnv(e.env.Derive())
	for _, expr := range expr.Exprs {
		var ret object.Object
		if ret = derivedEvaluator.expectEvalToAnyType(expr); object.IsError(ret) {
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
		ret = &object.Null{}
	}

	return ret
}
func (e *Evaluator) evalGtExpression(expr ast.GtExpression) object.Object {
	var left, right object.Object

	if left = e.expectEvalToType(expr.Left, object.NUMBER); object.IsError(left) {
		return left
	}
	if right = e.expectEvalToType(expr.Right, object.NUMBER); object.IsError(right) {
		return right
	}

	return &object.Boolean{Value: left.(*object.Number).Value > right.(*object.Number).Value}
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
	if _, ok := builtinFunctions[expr.Name]; !ok {
		panic("unknown built-in function: " + expr.Name)
	}
	args := map[string]object.Object{}
	for _, argName := range builtinFunctions[expr.Name].arguments {
		v, ok := e.env.Get(argName)
		if !ok {
			return &object.Error{Msg: fmt.Sprintf("`%s` parameter required", argName)}
		}
		args[argName] = v
	}
	return builtinFunctions[expr.Name].body(e.env, args)
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
		Fields: map[string]object.Object{},
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
		ret.Fields[field.(object.Hashable).Hash()] = fieldVal
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
	var s object.Object
	if s = e.expectEvalToAnyType(expr.Left); object.IsError(s) {
		return s
	}

	var rval object.Object
	var val object.Object
	if s.Type() == object.STRUCT {
		if rval = e.expectEvalToType(expr.Right, object.STRING); object.IsError(rval) {
			return rval
		}
		var ok bool
		val, ok = s.(*object.Struct).Fields[rval.(*object.String).Value]
		if !ok {
			return &object.Error{Msg: "field does not exist: " + rval.(*object.String).Value, Loc: expr.Right.Location()}
		}
	} else if s.Type() == object.MAP {
		if rval = e.expectEvalToHashableType(expr.Right); object.IsError(rval) {
			return rval
		}
		var ok bool
		val, ok = s.(*object.Map).Fields[rval.(object.Hashable).Hash()]
		if !ok {
			return &object.Error{Msg: "map item does not exist: " + rval.(object.Hashable).Hash()}
		}
	} else if s.Type() == object.MODULE {
		if rval = e.expectEvalToType(expr.Right, object.STRING); object.IsError(rval) {
			return rval
		}
		var ok bool
		val, ok = s.(*object.Module).Exports[rval.(*object.String).Value]
		if !ok {
			return &object.Error{Msg: "identifier is not exported: " + rval.(*object.String).Value, Loc: expr.Right.Location()}
		}
	} else if s.Type() == object.ARRAY {
		if rval = e.expectEvalToType(expr.Right, object.NUMBER); object.IsError(rval) {
			return rval
		}
		index := rval.(*object.Number).Value

		if index >= len(s.(*object.Array).Items) || index < 0 {
			return &object.Error{Msg: "index out of range: " + strconv.Itoa(index)}
		}
		val = s.(*object.Array).Items[index]
	} else if s.Type() == object.STRING {
		if rval = e.expectEvalToType(expr.Right, object.NUMBER); object.IsError(rval) {
			return rval
		}
		index := rval.(*object.Number).Value

		if index >= len(s.(*object.String).Value) || index < 0 {
			return &object.Error{Msg: "index out of range: " + strconv.Itoa(index)}
		}
		val = &object.String{Value: string(s.(*object.String).Value[index])}
	} else {
		return &object.Error{Msg: "field access operator is not supported on this type: " + s.Type().String(), Loc: expr.Left.Location()}
	}

	return val
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
	var ret object.Object = &object.Null{}

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
	}

	return ret
}
func (e *Evaluator) evalForExpression(expr ast.ForExpression) object.Object {
	var r object.Object
	var ret object.Object = &object.Null{}

	// todo: support maps
	if r = e.expectEvalToType(expr.Range, object.ARRAY); object.IsError(r) {
		return r
	}

	_, returnArrowExpression := expr.Body.(ast.ArrowExpression)
	arrowItems := &object.Array{}

	derivedEvaluator := NewWithEnv(e.env.Derive())
	for i, v := range r.(*object.Array).Items {
		if expr.Index != nil {
			derivedEvaluator.env.Set(expr.Index.Name, &object.Number{Value: i})
		}
		derivedEvaluator.env.Set(expr.Value.Name, v)

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
	return &object.Boolean{Value: !val.(*object.Boolean).Value}
}
func (e *Evaluator) evalEqTestExpression(expr ast.EqTestExpression) object.Object {
	var lval, rval object.Object
	if lval = e.expectEvalToAnyType(expr.Left); object.IsError(lval) {
		return lval
	}
	if rval = e.expectEvalToAnyType(expr.Right); object.IsError(rval) {
		return rval
	}
	if lval.Type() != rval.Type() {
		return &object.Boolean{Value: false}
	}
	if lval.Type() == object.NUMBER {
		return &object.Boolean{Value: lval.(*object.Number).Value == rval.(*object.Number).Value}
	} else if lval.Type() == object.STRING {
		return &object.Boolean{Value: lval.(*object.String).Value == rval.(*object.String).Value}
	} else if lval.Type() == object.BOOLEAN {
		return &object.Boolean{Value: lval.(*object.Boolean).Value == rval.(*object.Boolean).Value}
	} else if lval.Type() == object.NULL {
		return &object.Boolean{Value: true}
	}
	return &object.Error{Msg: "don't know how to compare values of type " + lval.Type().String()}
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
