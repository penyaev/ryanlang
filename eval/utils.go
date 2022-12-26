package eval

import (
	"ryanlang/ast"
	"ryanlang/object"
)

func (e *Evaluator) expectEvalToAnyType(expr ast.Expression) object.Object {
	var result object.Object
	if result = e.Eval(expr); object.IsError(result) {
		return object.WrapError(result, "evaluating "+expr.String(), expr.Location())
	}

	return result
}

func (e *Evaluator) expectEvalToType(expr ast.Expression, typ object.Type) object.Object {
	var result object.Object
	if result = e.expectEvalToAnyType(expr); object.IsError(result) {
		return result
	}
	if err := object.ExpectType(result, typ); err != nil {
		return object.WrapError(err, expr.String()+" evaluated to "+result.String(), expr.Location())
	}

	return result
}

func (e *Evaluator) expectEvalToHashableType(expr ast.Expression) object.Object {
	var result object.Object
	if result = e.expectEvalToAnyType(expr); object.IsError(result) {
		return result
	}
	if !object.IsHashable(result) {
		return &object.Error{Msg: "expected hashable type, got: " + result.Type().String()}
	}

	return result
}
