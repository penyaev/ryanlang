package parser

import (
	"ryanlang/ast"
	"ryanlang/lexer"
)

type Parser struct {
	l   *lexer.Lexer
	cur *lexer.Token

	prefixFunctions map[lexer.TokenKind]prefixParseFunction
	infixFunctions  map[lexer.TokenKind]infixParseFunction
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.nextToken()
	p.prefixFunctions = map[lexer.TokenKind]prefixParseFunction{
		lexer.TokenTypeNumber:         p.parseNumber,
		lexer.TokenTypeString:         p.parseString,
		lexer.TokenTypeNull:           p.parseNull,
		lexer.TokenTypeLet:            p.parseLet,
		lexer.TokenTypeIdentifier:     p.parseIdentifier,
		lexer.TokenTypeImport:         p.parseImport,
		lexer.TokenTypeIf:             p.parseIf,
		lexer.TokenTypeLBracket:       p.parseGroup,
		lexer.TokenTypeFunc:           p.parseFunc,
		lexer.TokenTypeReturn:         p.parseReturn,
		lexer.TokenTypeContinue:       p.parseContinue,
		lexer.TokenTypeBreak:          p.parseBreak,
		lexer.TokenTypeStruct:         p.parseStruct,
		lexer.TokenTypeMap:            p.parseMap,
		lexer.TokenTypeExports:        p.parseExports,
		lexer.TokenTypeWhile:          p.parseWhile,
		lexer.TokenTypeFor:            p.parseFor,
		lexer.TokenTypeLSquareBracket: p.parseArray,
		lexer.TokenTypeBang:           p.parseNegation,
		lexer.TokenTypeMinus:          p.parsePrefixMinus,
	}

	p.infixFunctions = map[lexer.TokenKind]infixParseFunction{
		lexer.TokenTypePlus:        p.parsePlus,
		lexer.TokenTypePlusAssign:  p.parsePlusAssign,
		lexer.TokenTypePlusPlus:    p.parsePlusPlus,
		lexer.TokenTypeMinus:       p.parseMinus,
		lexer.TokenTypeMinusAssign: p.parseMinusAssign,
		lexer.TokenTypeMinusMinus:  p.parseMinusMinus,
		lexer.TokenTypeLogicalAnd:  p.parseLogicalAnd,
		lexer.TokenTypeLogicalOr:   p.parseLogicalOr,
		lexer.TokenTypeAsterisk:    p.parseMult,
		lexer.TokenTypeDiv:         p.parseDiv,
		lexer.TokenTypeLBracket:    p.parseCall,
		lexer.TokenTypeGt:          p.parseGt,
		lexer.TokenTypeLt:          p.parseLt,
		lexer.TokenTypeGte:         p.parseGte,
		lexer.TokenTypeLte:         p.parseLte,
		lexer.TokenTypeEqTest:      p.parseEqTest,
		lexer.TokenTypeNeqTest:     p.parseNeqTest,
		lexer.TokenTypeAssign:      p.parseAssign,
		lexer.TokenTypeDot:         p.parseFieldAccess,
		lexer.TokenTypeMod:         p.parseMod,
		lexer.TokenTypeComma:       p.parseTuple,
	}
	return p
}

func (p *Parser) nextToken() *lexer.Token {
	var err error
	for {
		p.cur, err = p.l.Next()
		if err != nil /*&& !errors.Is(err, lexer.ErrEof)*/ {
			panic(err)
		}
		if p.cur.Kind != lexer.TokenTypeComment {
			break
		}
	}
	return p.cur
}
func (p *Parser) location() string {
	return p.cur.Location.String()
}
func (p *Parser) consume(kind lexer.TokenKind) *lexer.Token {
	cur := p.cur
	if kind != cur.Kind {
		panic(p.location() + ": expected token: " + kind.String() + ", got: " + cur.Kind.String())
	}
	p.nextToken()
	return cur
}
func (p *Parser) curPrecedence() int {
	pr, ok := precedences[p.cur.Kind]
	if !ok {
		panic(p.location() + ": unknown precedence for " + p.cur.Kind.String())
	}
	return pr
}
func (p *Parser) readExpression(precedence int) ast.Expression {
	prefixParseFn, ok := p.prefixFunctions[p.cur.Kind]
	if !ok {
		panic(p.location() + ": no prefix function for " + p.cur.Kind.String())
	}
	expr := prefixParseFn()

	for p.cur.Kind != lexer.TokenTypeSemicolon && precedence < p.curPrecedence() {
		infixParseFn, ok := p.infixFunctions[p.cur.Kind]
		if !ok {
			panic("unexpected infix operator: " + p.cur.Kind.String())
		}

		expr = infixParseFn(expr)
	}

	return expr
}
func (p *Parser) Eof() bool {
	return p.cur.Kind == lexer.TokenTypeEof
}
func (p *Parser) ReadStatement() ast.Expression {
	expr := p.readExpression(precedenceLowest)
	p.consume(lexer.TokenTypeSemicolon)
	return expr
}
func (p *Parser) ReadModule(name string) ast.Expression {
	ret := ast.Module{
		Name: name,
	}
	for !p.Eof() {
		st := p.ReadStatement()
		if _, ok := st.(ast.Exports); ok && len(ret.Exprs) > 0 {
			panic("`exports` must come first statement in a module")
		}
		ret.Exprs = append(ret.Exprs, st)
	}
	return ret
}
