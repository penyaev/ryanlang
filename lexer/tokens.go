package lexer

import (
	"fmt"
	"sort"
)

type TokenKind int

const (
	_ TokenKind = iota
	TokenTypeLet
	TokenTypeIdentifier
	TokenTypeAssign
	TokenTypeEqTest
	TokenTypeNeqTest
	TokenTypeNumber
	TokenTypeLBracket
	TokenTypeRBracket
	TokenTypeLSquareBracket
	TokenTypeRSquareBracket
	TokenTypeComma
	TokenTypeIn
	TokenTypeUntil // todo: unused
	TokenTypeSemicolon
	TokenTypeIf
	TokenTypeElse
	TokenTypeGt
	TokenTypeGte
	TokenTypeLt
	TokenTypeLte
	TokenTypeLBrace
	TokenTypeRBrace
	TokenTypeReturn
	TokenTypePlus
	TokenTypePlusPlus
	TokenTypePlusAssign
	TokenTypeMinus
	TokenTypeMinusMinus
	TokenTypeMinusAssign
	TokenTypeAsterisk
	TokenTypeDiv
	TokenTypeFunc
	TokenTypeEof
	TokenTypeComment
	TokenTypeString
	TokenTypeArrow
	TokenTypeStruct
	TokenTypeColon
	TokenTypeDot
	TokenTypeWhile
	TokenTypeBang
	TokenTypeLogicalAnd
	TokenTypeLogicalOr
	TokenTypeContinue
	TokenTypeBreak
	TokenTypeFor
	TokenTypeMod
	TokenTypeExports
	TokenTypeImport
	TokenTypeMap
)

func (tk TokenKind) String() string {
	switch tk {
	case TokenTypeAssign:
		return "="
	case TokenTypeEqTest:
		return "=="
	case TokenTypeNeqTest:
		return "!="
	case TokenTypeSemicolon:
		return ";"
	case TokenTypePlus:
		return "+"
	case TokenTypePlusPlus:
		return "++"
	case TokenTypePlusAssign:
		return "+="
	case TokenTypeMinus:
		return "-"
	case TokenTypeMinusMinus:
		return "--"
	case TokenTypeMinusAssign:
		return "-="
	case TokenTypeLBracket:
		return "("
	case TokenTypeComma:
		return ","
	case TokenTypeRBracket:
		return ")"
	case TokenTypeAsterisk:
		return "*"
	case TokenTypeDiv:
		return "/"
	case TokenTypeIf:
		return "if"
	case TokenTypeElse:
		return "else"
	case TokenTypeGt:
		return ">"
	case TokenTypeLt:
		return "<"
	case TokenTypeGte:
		return ">="
	case TokenTypeLte:
		return "<="
	case TokenTypeLBrace:
		return "{"
	case TokenTypeRBrace:
		return "}"
	case TokenTypeFunc:
		return "func"
	case TokenTypeNumber:
		return "<number>"
	case TokenTypeReturn:
		return "return"
	case TokenTypeIdentifier:
		return "<identifier>"
	case TokenTypeEof:
		return "<eof>"
	case TokenTypeComment:
		return "<comment>"
	case TokenTypeString:
		return "<string>"
	case TokenTypeArrow:
		return "=>"
	case TokenTypeStruct:
		return "<struct>"
	case TokenTypeColon:
		return ":"
	case TokenTypeDot:
		return "."
	case TokenTypeWhile:
		return "while"
	case TokenTypeLSquareBracket:
		return "["
	case TokenTypeRSquareBracket:
		return "]"
	case TokenTypeLet:
		return "let"
	case TokenTypeBang:
		return "!"
	case TokenTypeLogicalAnd:
		return "&&"
	case TokenTypeLogicalOr:
		return "||"
	case TokenTypeContinue:
		return "continue"
	case TokenTypeBreak:
		return "break"
	case TokenTypeFor:
		return "for"
	case TokenTypeIn:
		return "in"
	case TokenTypeUntil:
		return "until"
	case TokenTypeMod:
		return "%"
	case TokenTypeExports:
		return "exports"
	case TokenTypeMap:
		return "map"
	case TokenTypeImport:
		return "import"
	default:
		return fmt.Sprintf("[%d]", tk)
	}
}

var keywords = map[string]TokenKind{
	"if":       TokenTypeIf,
	"else":     TokenTypeElse,
	"let":      TokenTypeLet,
	"return":   TokenTypeReturn,
	"func":     TokenTypeFunc,
	"struct":   TokenTypeStruct,
	"while":    TokenTypeWhile,
	"continue": TokenTypeContinue,
	"break":    TokenTypeBreak,
	"for":      TokenTypeFor,
	"in":       TokenTypeIn,
	"until":    TokenTypeUntil,
	"exports":  TokenTypeExports,
	"import":   TokenTypeImport,
	"map":      TokenTypeMap,
}
var tokens = []struct {
	literal string
	kind    TokenKind
}{
	{"+", TokenTypePlus},
	{"++", TokenTypePlusPlus},
	{"+=", TokenTypePlusAssign},
	{"-", TokenTypeMinus},
	{"--", TokenTypeMinusMinus},
	{"-=", TokenTypeMinusAssign},
	{">", TokenTypeGt},
	{"<", TokenTypeLt},
	{">=", TokenTypeGte},
	{"<=", TokenTypeLte},
	{"{", TokenTypeLBrace},
	{"}", TokenTypeRBrace},
	{"(", TokenTypeLBracket},
	{")", TokenTypeRBracket},
	{"[", TokenTypeLSquareBracket},
	{"]", TokenTypeRSquareBracket},
	{",", TokenTypeComma},
	{";", TokenTypeSemicolon},
	{"=", TokenTypeAssign},
	{"==", TokenTypeEqTest},
	{"!=", TokenTypeNeqTest},
	{"*", TokenTypeAsterisk},
	{"/", TokenTypeDiv},
	{"=>", TokenTypeArrow},
	{":", TokenTypeColon},
	{".", TokenTypeDot},
	{"!", TokenTypeBang},
	{"&&", TokenTypeLogicalAnd},
	{"||", TokenTypeLogicalOr},
	{"%", TokenTypeMod},
}

func init() {
	sort.Slice(tokens, func(i, j int) bool {
		return len(tokens[i].literal) > len(tokens[j].literal)
	})
}

type Token struct {
	Kind     TokenKind
	Literal  string
	Location *Location
}

type Location struct {
	File   string
	Line   int
	Column int
}

func (l *Location) Clone() *Location {
	return &Location{
		File:   l.File,
		Line:   l.Line,
		Column: l.Column,
	}
}

func (l *Location) String() string {
	if l == nil {
		return "(unknown location)"
	}
	return fmt.Sprintf("%s:%d:%d", l.File, l.Line, l.Column)
}
