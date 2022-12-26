package parser

import (
	"reflect"
	"ryanlang/ast"
	"ryanlang/lexer"
)

type infixParseFunction func(ast.Expression) ast.Expression

func (p *Parser) parsePlus(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypePlus)
	return ast.PlusExpression{
		Left:  left,
		Right: p.readExpression(precedencePlusMinus),
	}
}
func (p *Parser) parsePlusPlus(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypePlusPlus)
	switch left.(type) {
	case ast.Identifier:
		identifier := left.(ast.Identifier)
		return ast.AssignExpression{
			Identifier: identifier,
			Value: ast.PlusExpression{
				Left:  identifier,
				Right: ast.NumberExpression{Value: 1},
			},
		}
	case ast.FieldAccessExpression:
		fieldAccess := left.(ast.FieldAccessExpression)
		return ast.FieldAssignExpression{
			FieldAccess: fieldAccess,
			Value: ast.PlusExpression{
				Left:  fieldAccess,
				Right: ast.NumberExpression{Value: 1},
			},
		}
	default:
		panic("identifier or a dot-expression expected on the left side of the incr operator, got: " + reflect.TypeOf(left).String())
	}
}
func (p *Parser) parsePlusAssign(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypePlusAssign)
	switch left.(type) {
	case ast.Identifier:
		identifier := left.(ast.Identifier)
		return ast.AssignExpression{
			Identifier: identifier,
			Value: ast.PlusExpression{
				Left:  identifier,
				Right: p.readExpression(precedenceAssign),
			},
		}
	case ast.FieldAccessExpression:
		fieldAccess := left.(ast.FieldAccessExpression)
		return ast.FieldAssignExpression{
			FieldAccess: fieldAccess,
			Value: ast.PlusExpression{
				Left:  fieldAccess,
				Right: p.readExpression(precedenceAssign),
			},
		}
	default:
		panic("identifier or a dot-expression expected on the left side of the plus-assignment operator, got: " + reflect.TypeOf(left).String())
	}
}

func (p *Parser) parseMinus(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeMinus)
	return ast.MinusExpression{
		Left:  left,
		Right: p.readExpression(precedencePlusMinus),
	}
}
func (p *Parser) parseMinusMinus(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeMinusMinus)
	switch left.(type) {
	case ast.Identifier:
		identifier := left.(ast.Identifier)
		return ast.AssignExpression{
			Identifier: identifier,
			Value: ast.MinusExpression{
				Left:  identifier,
				Right: ast.NumberExpression{Value: 1},
			},
		}
	case ast.FieldAccessExpression:
		fieldAccess := left.(ast.FieldAccessExpression)
		return ast.FieldAssignExpression{
			FieldAccess: fieldAccess,
			Value: ast.MinusExpression{
				Left:  fieldAccess,
				Right: ast.NumberExpression{Value: 1},
			},
		}
	default:
		panic("identifier or a dot-expression expected on the left side of the decr operator, got: " + reflect.TypeOf(left).String())
	}
}
func (p *Parser) parseMinusAssign(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeMinusAssign)
	switch left.(type) {
	case ast.Identifier:
		identifier := left.(ast.Identifier)
		return ast.AssignExpression{
			Identifier: identifier,
			Value: ast.MinusExpression{
				Left:  identifier,
				Right: p.readExpression(precedenceAssign),
			},
		}
	case ast.FieldAccessExpression:
		fieldAccess := left.(ast.FieldAccessExpression)
		return ast.FieldAssignExpression{
			FieldAccess: fieldAccess,
			Value: ast.MinusExpression{
				Left:  fieldAccess,
				Right: p.readExpression(precedenceAssign),
			},
		}
	default:
		panic("identifier or a dot-expression expected on the left side of the minus-assignment operator, got: " + reflect.TypeOf(left).String())
	}
}

func (p *Parser) parseMod(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeMod)
	return ast.ModExpression{
		Left:  left,
		Right: p.readExpression(precedenceMultDiv),
	}
}
func (p *Parser) parseLogicalAnd(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeLogicalAnd)
	return ast.LogicalAndExpression{
		Left:  left,
		Right: p.readExpression(precedenceLogicalAnd),
	}
}
func (p *Parser) parseLogicalOr(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeLogicalOr)
	return ast.LogicalOrExpression{
		Left:  left,
		Right: p.readExpression(precedenceLogicalOr),
	}
}

func (p *Parser) parseMult(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeAsterisk)
	return ast.MultExpression{
		Left:  left,
		Right: p.readExpression(precedenceMultDiv),
	}
}

func (p *Parser) parseDiv(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeDiv)
	return ast.DivExpression{
		Left:  left,
		Right: p.readExpression(precedenceMultDiv),
	}
}

func (p *Parser) parseGt(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeGt)
	return ast.GtExpression{
		Left:  left,
		Right: p.readExpression(precedenceGtLt),
	}
}
func (p *Parser) parseLt(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeLt)
	return ast.LtExpression{
		Left:  left,
		Right: p.readExpression(precedenceGtLt),
	}
}
func (p *Parser) parseGte(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeGte)
	return ast.GteExpression{
		Left:  left,
		Right: p.readExpression(precedenceGtLt),
	}
}
func (p *Parser) parseLte(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeLte)
	return ast.LteExpression{
		Left:  left,
		Right: p.readExpression(precedenceGtLt),
	}
}
func (p *Parser) parseTuple(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeComma)
	var exp []ast.Expression
	if _, ok := left.(ast.TupleExpression); ok {
		exp = left.(ast.TupleExpression).Exprs
	} else {
		exp = append(exp, left)
	}

	return ast.TupleExpression{
		Exprs: append(exp, p.readExpression(precedenceComma)),
	}
}

func (p *Parser) readCommaSeparatedExpressions(endToken lexer.TokenKind) []ast.Expression {
	var ret []ast.Expression
	expectingExpression := false
	for p.cur.Kind != endToken {
		ret = append(ret, p.readExpression(precedenceComma))
		expectingExpression = false
		if p.cur.Kind == lexer.TokenTypeComma {
			p.consume(lexer.TokenTypeComma)
			expectingExpression = true
		} else if p.cur.Kind != endToken {
			panic(p.location() + ": expected: comma or " + endToken.String())
		}
	}
	if expectingExpression {
		panic(p.location() + ": expression expected after comma")
	}
	return ret
}
func (p *Parser) parseCall(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeLBracket)
	args := p.readCommaSeparatedExpressions(lexer.TokenTypeRBracket)
	p.consume(lexer.TokenTypeRBracket)
	return ast.CallExpression{
		Callee:    left,
		Arguments: args,
	}
}

func (p *Parser) parseEqTest(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeEqTest)
	return ast.EqTestExpression{
		Left:  left,
		Right: p.readExpression(precedenceEqTest),
	}
}
func (p *Parser) parseNeqTest(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeNeqTest)
	return ast.NegationExpression{
		Expr: ast.EqTestExpression{
			Left:  left,
			Right: p.readExpression(precedenceEqTest),
		},
	}
}

func (p *Parser) parseAssign(left ast.Expression) ast.Expression {
	switch left.(type) {
	case ast.Identifier:
		identifier := left.(ast.Identifier)
		p.consume(lexer.TokenTypeAssign)
		return ast.AssignExpression{
			Identifier: identifier,
			Value:      p.readExpression(precedenceAssign),
		}
	case ast.FieldAccessExpression:
		fieldAccess := left.(ast.FieldAccessExpression)
		p.consume(lexer.TokenTypeAssign)
		return ast.FieldAssignExpression{
			FieldAccess: fieldAccess,
			Value:       p.readExpression(precedenceAssign),
		}
	case ast.TupleExpression:
		tuple := left.(ast.TupleExpression)
		p.consume(lexer.TokenTypeAssign)
		return ast.TupleAssignExpression{
			Tuple: tuple,
			Value: p.readExpression(precedenceAssign),
		}
	default:
		panic("identifier or a dot-expression expected on the left side of the assignment operator, got: " + reflect.TypeOf(left).String())
	}
}
func (p *Parser) parseFieldAccess(left ast.Expression) ast.Expression {
	p.consume(lexer.TokenTypeDot)

	right := p.readExpression(precedenceFieldAccess)
	switch right.(type) {
	case ast.Identifier:
		right = ast.String{Value: right.(ast.Identifier).Name, Loc: right.(ast.Identifier).Location()}
	}
	return ast.FieldAccessExpression{
		Left:  left,
		Right: right,
	}
}
