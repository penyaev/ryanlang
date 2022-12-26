package eval

import (
	"reflect"
	"ryanlang/ast"
	"ryanlang/object"
)

type Evaluator struct {
	env *object.Environment
}

func New() *Evaluator {
	return &Evaluator{env: newBuiltinEnvironment()}
}

func NewWithEnv(env *object.Environment) *Evaluator {
	return &Evaluator{env: env}
}

func (e *Evaluator) Eval(expr ast.Expression) object.Object {
	switch expr.(type) {
	case ast.PlusExpression:
		return e.evalPlusExpression(expr.(ast.PlusExpression))
	case ast.MinusExpression:
		return e.evalMinusExpression(expr.(ast.MinusExpression))
	case ast.ModExpression:
		return e.evalModExpression(expr.(ast.ModExpression))
	case ast.MultExpression:
		return e.evalMultExpression(expr.(ast.MultExpression))
	case ast.DivExpression:
		return e.evalDivExpression(expr.(ast.DivExpression))
	case ast.LogicalAndExpression:
		return e.evalLogicalAndExpression(expr.(ast.LogicalAndExpression))
	case ast.LogicalOrExpression:
		return e.evalLogicalOrExpression(expr.(ast.LogicalOrExpression))
	case ast.LetExpression:
		return e.evalLetExpression(expr.(ast.LetExpression))
	case ast.Identifier:
		return e.evalIdentifier(expr.(ast.Identifier))
	case ast.String:
		return e.evalString(expr.(ast.String))
	case ast.Null:
		return e.evalNull(expr.(ast.Null))
	case ast.NumberExpression:
		return e.evalNumber(expr.(ast.NumberExpression))
	case ast.AssignExpression:
		return e.evalAssignExpression(expr.(ast.AssignExpression))
	case ast.FieldAssignExpression:
		return e.evalFieldAssignExpression(expr.(ast.FieldAssignExpression))
	case ast.FuncExpression:
		return e.evalFuncExpression(expr.(ast.FuncExpression))
	case ast.CallExpression:
		return e.evalCallExpression(expr.(ast.CallExpression))
	case ast.ReturnExpression:
		return e.evalReturnExpression(expr.(ast.ReturnExpression))
	case ast.ContinueExpression:
		return e.evalContinueExpression(expr.(ast.ContinueExpression))
	case ast.BreakExpression:
		return e.evalBreakExpression(expr.(ast.BreakExpression))
	case ast.BlockExpression:
		return e.evalBlockExpression(expr.(ast.BlockExpression))
	case ast.GroupExpression:
		return e.evalGroupExpression(expr.(ast.GroupExpression))
	case ast.IfExpression:
		return e.evalIfExpression(expr.(ast.IfExpression))
	case ast.GtExpression:
		return e.evalGtExpression(expr.(ast.GtExpression))
	case ast.LtExpression:
		return e.evalLtExpression(expr.(ast.LtExpression))
	case ast.GteExpression:
		return e.evalGteExpression(expr.(ast.GteExpression))
	case ast.LteExpression:
		return e.evalLteExpression(expr.(ast.LteExpression))
	case ast.BuiltinFunction:
		return e.evalBuiltinFunction(expr.(ast.BuiltinFunction))
	case ast.StructExpression:
		return e.evalStructExpression(expr.(ast.StructExpression))
	case ast.MapExpression:
		return e.evalMapExpression(expr.(ast.MapExpression))
	case ast.Exports:
		return e.evalExports(expr.(ast.Exports))
	case ast.FieldAccessExpression:
		return e.evalFieldAccessExpression(expr.(ast.FieldAccessExpression))
	case ast.WhileExpression:
		return e.evalWhileExpression(expr.(ast.WhileExpression))
	case ast.ForExpression:
		return e.evalForExpression(expr.(ast.ForExpression))
	case ast.ArrayExpression:
		return e.evalArrayExpression(expr.(ast.ArrayExpression))
	case ast.NegationExpression:
		return e.evalNegationExpression(expr.(ast.NegationExpression))
	case ast.EqTestExpression:
		return e.evalEqTestExpression(expr.(ast.EqTestExpression))
	case ast.PrefixMinusExpression:
		return e.evalPrefixMinusExpression(expr.(ast.PrefixMinusExpression))
	case ast.Module:
		return e.evalModule(expr.(ast.Module))
	case ast.Import:
		return e.evalImport(expr.(ast.Import))
	case ast.ArrowExpression:
		return e.evalArrowExpression(expr.(ast.ArrowExpression))
	case ast.TupleExpression:
		return e.evalTupleExpression(expr.(ast.TupleExpression))
	case ast.TupleAssignExpression:
		return e.evalTupleAssignExpression(expr.(ast.TupleAssignExpression))
	}

	panic("missing eval implementation for " + reflect.TypeOf(expr).String())
}
