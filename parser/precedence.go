package parser

import "ryanlang/lexer"

// see for inspiration: https://en.cppreference.com/w/c/language/operator_precedence
const (
	precedenceLowest = iota
	precedenceAssign
	precedenceComma
	precedenceLogicalOr
	precedenceLogicalAnd
	precedenceEqTest
	precedenceGtLt
	precedencePlusMinus
	precedenceBang
	precedenceMultDiv
	precedencePrefixMinus
	precedenceIncrDecr
	precedenceCall
	precedenceFieldAccess
	precedenceHighest
)

var precedences = map[lexer.TokenKind]int{
	lexer.TokenTypePlus:           precedencePlusMinus,
	lexer.TokenTypeMinus:          precedencePlusMinus,
	lexer.TokenTypeAsterisk:       precedenceMultDiv,
	lexer.TokenTypeDiv:            precedenceMultDiv,
	lexer.TokenTypeLBracket:       precedenceCall,
	lexer.TokenTypeComma:          precedenceComma,
	lexer.TokenTypeRBracket:       precedenceLowest,
	lexer.TokenTypeGt:             precedenceGtLt,
	lexer.TokenTypeLt:             precedenceGtLt,
	lexer.TokenTypeGte:            precedenceGtLt,
	lexer.TokenTypeLte:            precedenceGtLt,
	lexer.TokenTypeLBrace:         precedenceLowest,
	lexer.TokenTypeRBrace:         precedenceLowest,
	lexer.TokenTypeAssign:         precedenceAssign,
	lexer.TokenTypeNumber:         precedenceLowest,
	lexer.TokenTypeEqTest:         precedenceEqTest,
	lexer.TokenTypeNeqTest:        precedenceEqTest,
	lexer.TokenTypeEof:            precedenceLowest,
	lexer.TokenTypeElse:           precedenceLowest,
	lexer.TokenTypeArrow:          precedenceLowest,
	lexer.TokenTypeIdentifier:     precedenceLowest,
	lexer.TokenTypeDot:            precedenceFieldAccess,
	lexer.TokenTypeRSquareBracket: precedenceLowest,
	lexer.TokenTypeLogicalAnd:     precedenceLogicalAnd,
	lexer.TokenTypeLogicalOr:      precedenceLogicalOr,
	lexer.TokenTypeMod:            precedenceMultDiv,
	lexer.TokenTypePlusAssign:     precedenceAssign,
	lexer.TokenTypeMinusAssign:    precedenceAssign,
	lexer.TokenTypePlusPlus:       precedenceIncrDecr,
	lexer.TokenTypeMinusMinus:     precedenceIncrDecr,
	lexer.TokenTypeColon:          precedenceLowest,
}
