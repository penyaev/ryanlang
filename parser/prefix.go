package parser

import (
	"ryanlang/ast"
	"ryanlang/lexer"
	"strconv"
)

type prefixParseFunction func() ast.Expression

func (p *Parser) parseNumber() ast.Expression {
	number, err := strconv.Atoi(p.consume(lexer.TokenTypeNumber).Literal)
	if err != nil {
		panic(err)
	}
	return ast.NumberExpression{
		Value: number,
	}
}
func (p *Parser) parseString() ast.Expression {
	loc := p.cur.Location
	return ast.String{
		Loc:   loc,
		Value: p.consume(lexer.TokenTypeString).Literal,
	}
}
func (p *Parser) parseNull() ast.Expression {
	p.consume(lexer.TokenTypeNull)
	return ast.Null{}
}
func (p *Parser) parseIdentifier() ast.Expression {
	loc := p.cur.Location
	return ast.Identifier{
		Loc:  loc,
		Name: p.consume(lexer.TokenTypeIdentifier).Literal,
	}
}
func (p *Parser) parseImport() ast.Expression {
	loc := p.cur.Location
	p.consume(lexer.TokenTypeImport)
	p.consume(lexer.TokenTypeLBracket)
	expr := p.readExpression(precedenceLowest)
	p.consume(lexer.TokenTypeRBracket)
	return ast.Import{
		Loc:    loc,
		Module: expr,
	}
}
func (p *Parser) parseLet() ast.Expression {
	loc := p.cur.Location
	p.consume(lexer.TokenTypeLet)
	ids := []ast.Identifier{}
	ids = append(ids, p.parseIdentifier().(ast.Identifier))
	for p.cur.Kind == lexer.TokenTypeComma {
		p.consume(lexer.TokenTypeComma)
		ids = append(ids, p.parseIdentifier().(ast.Identifier))
	}
	p.consume(lexer.TokenTypeAssign)
	initialization := p.readExpression(precedenceAssign)
	return ast.LetExpression{
		Identifiers:    ids,
		Initialization: initialization,
		Loc:            loc,
	}
}
func (p *Parser) parseIf() ast.Expression {
	loc := p.cur.Location
	p.consume(lexer.TokenTypeIf)
	result := ast.IfExpression{
		Condition: p.readExpression(precedenceLowest),
		Loc:       loc,
	}

	if p.cur.Kind == lexer.TokenTypeLBrace {
		result.Then = p.readBlockExpression()
	} else if p.cur.Kind == lexer.TokenTypeArrow {
		result.Then = p.readArrowExpression()
	} else {
		panic("`then` expected as an arrow expression or a statement block")
	}

	if p.cur.Kind == lexer.TokenTypeElse {
		p.consume(lexer.TokenTypeElse)

		if p.cur.Kind == lexer.TokenTypeLBrace {
			result.Else = p.readBlockExpression()
		} else if p.cur.Kind == lexer.TokenTypeArrow {
			result.Else = p.readArrowExpression()
		} else if p.cur.Kind == lexer.TokenTypeIf { // chained if
			result.Else = p.parseIf()
		} else {
			panic("`else` expected as an arrow expression, if or a statement block")
		}
	}

	return result
}
func (p *Parser) parseGroup() ast.Expression {
	loc := p.cur.Location
	p.consume(lexer.TokenTypeLBracket)
	result := p.readExpression(precedenceLowest)
	p.consume(lexer.TokenTypeRBracket)

	return ast.GroupExpression{Expr: result, Loc: loc}
}
func (p *Parser) readBlockExpression() ast.BlockExpression {
	result := ast.BlockExpression{
		Loc: p.cur.Location,
	}
	p.consume(lexer.TokenTypeLBrace)
	for p.cur.Kind != lexer.TokenTypeRBrace {
		result.Exprs = append(result.Exprs, p.ReadStatement())
	}
	p.consume(lexer.TokenTypeRBrace)
	return result
}
func (p *Parser) readArrowExpression() ast.Expression {
	loc := p.cur.Location
	p.consume(lexer.TokenTypeArrow)
	return ast.ArrowExpression{Expr: p.readExpression(precedenceLowest), Loc: loc}
}
func (p *Parser) parseFunc() ast.Expression {
	result := ast.FuncExpression{}

	p.consume(lexer.TokenTypeFunc)

	if p.cur.Kind == lexer.TokenTypeLBracket {
		p.consume(lexer.TokenTypeLBracket)
		expectingArgument := false
		for p.cur.Kind != lexer.TokenTypeRBracket {
			result.Arguments = append(result.Arguments, p.parseIdentifier().(ast.Identifier))
			expectingArgument = false
			if p.cur.Kind == lexer.TokenTypeComma {
				p.consume(lexer.TokenTypeComma)
				expectingArgument = true
			}
		}
		if expectingArgument {
			panic("argument expected after comma")
		}
		p.consume(lexer.TokenTypeRBracket)
	}

	if p.cur.Kind == lexer.TokenTypeLBrace {
		result.Body = p.readBlockExpression()
	} else if p.cur.Kind == lexer.TokenTypeArrow {
		result.Body = ast.ReturnExpression{Expr: p.readArrowExpression()}
	} else {
		panic(p.cur.Location.String() + ": missing function body")
	}

	return result
}
func (p *Parser) parseReturn() ast.Expression {
	p.consume(lexer.TokenTypeReturn)
	return ast.ReturnExpression{
		Expr: p.readExpression(precedenceLowest),
	}
}
func (p *Parser) parseContinue() ast.Expression {
	p.consume(lexer.TokenTypeContinue)
	return ast.ContinueExpression{}
}
func (p *Parser) parseBreak() ast.Expression {
	p.consume(lexer.TokenTypeBreak)
	return ast.BreakExpression{}
}
func (p *Parser) parseStruct() ast.Expression {
	result := ast.StructExpression{
		Fields: map[string]ast.Expression{},
	}
	p.consume(lexer.TokenTypeStruct)
	p.consume(lexer.TokenTypeLBrace)
	for p.cur.Kind != lexer.TokenTypeRBrace {
		id := p.parseIdentifier()
		p.consume(lexer.TokenTypeColon)
		if _, ok := result.Fields[id.(ast.Identifier).Name]; ok {
			panic("duplicate struct field: " + id.(ast.Identifier).Name)
		}
		result.Fields[id.(ast.Identifier).Name] = p.readExpression(precedenceLowest)
		p.consume(lexer.TokenTypeSemicolon)
	}
	p.consume(lexer.TokenTypeRBrace)

	return result
}
func (p *Parser) parseMap() ast.Expression {
	result := ast.MapExpression{}
	p.consume(lexer.TokenTypeMap)
	p.consume(lexer.TokenTypeLBrace)
	for p.cur.Kind != lexer.TokenTypeRBrace {
		id := p.readExpression(precedenceLowest)
		p.consume(lexer.TokenTypeColon)
		result.Fields = append(result.Fields, ast.MapField{
			Key:   id,
			Value: p.readExpression(precedenceLowest),
		})
		p.consume(lexer.TokenTypeSemicolon)
	}
	p.consume(lexer.TokenTypeRBrace)

	return result
}
func (p *Parser) parseExports() ast.Expression {
	result := ast.Exports{
		Fields: map[string]ast.Expression{},
	}
	p.consume(lexer.TokenTypeExports)
	p.consume(lexer.TokenTypeLBrace)
	for p.cur.Kind != lexer.TokenTypeRBrace {
		id := p.parseIdentifier()
		if _, ok := result.Fields[id.(ast.Identifier).Name]; ok {
			panic("duplicate exports field: " + id.(ast.Identifier).Name)
		}
		var init ast.Expression = nil
		if p.cur.Kind == lexer.TokenTypeColon {
			p.consume(lexer.TokenTypeColon)
			init = p.readExpression(precedenceLowest)
		}
		result.Fields[id.(ast.Identifier).Name] = init
		p.consume(lexer.TokenTypeSemicolon)
	}
	p.consume(lexer.TokenTypeRBrace)

	return result
}
func (p *Parser) parseWhile() ast.Expression {
	p.consume(lexer.TokenTypeWhile)
	result := ast.WhileExpression{
		Condition: p.readExpression(precedenceLowest),
	}

	if p.cur.Kind == lexer.TokenTypeLBrace {
		result.Body = p.readBlockExpression()
	} else if p.cur.Kind == lexer.TokenTypeArrow {
		result.Body = p.readArrowExpression()
	} else {
		panic("missing while loop body")
	}

	return result
}
func (p *Parser) parseFor() ast.Expression {
	loc := p.cur.Location
	p.consume(lexer.TokenTypeFor)
	result := ast.ForExpression{
		Loc: loc,
	}

	result.Value = p.parseIdentifier().(ast.Identifier)
	if p.cur.Kind == lexer.TokenTypeComma {
		p.consume(lexer.TokenTypeComma)
		index := result.Value
		result.Index = &index
		result.Value = p.parseIdentifier().(ast.Identifier)
	}

	p.consume(lexer.TokenTypeIn)

	result.Range = p.readExpression(precedenceLowest)

	if p.cur.Kind == lexer.TokenTypeLBrace {
		result.Body = p.readBlockExpression()
	} else if p.cur.Kind == lexer.TokenTypeArrow {
		result.Body = p.readArrowExpression()
	} else {
		panic("missing for loop body")
	}

	return result
}
func (p *Parser) parseArray() ast.Expression {
	p.consume(lexer.TokenTypeLSquareBracket)
	items := p.readCommaSeparatedExpressions(lexer.TokenTypeRSquareBracket)
	p.consume(lexer.TokenTypeRSquareBracket)
	return ast.ArrayExpression{
		Items: items,
	}
}
func (p *Parser) parseNegation() ast.Expression {
	p.consume(lexer.TokenTypeBang)
	return ast.NegationExpression{
		Expr: p.readExpression(precedenceBang),
	}
}
func (p *Parser) parsePrefixMinus() ast.Expression {
	p.consume(lexer.TokenTypeMinus)
	return ast.PrefixMinusExpression{
		Expr: p.readExpression(precedencePrefixMinus),
	}
}
