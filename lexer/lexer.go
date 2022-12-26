package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var ErrEof = fmt.Errorf("eof")
var ErrUnknownByte = fmt.Errorf("unknown byte")
var ErrUnterminatedCommentBlock = fmt.Errorf("unterminated comment block")
var ErrUnterminatedString = fmt.Errorf("unterminated string")
var ErrInvalidNumberFormat = fmt.Errorf("invalid number format")

type Lexer struct {
	r       *bufio.Reader
	eof     bool
	buf     []rune
	loc     *Location
	prevLoc *Location
}

// todo: check assignments to values passed as parameters
// e.g.: func(x) { x = 10; }
// will it overwrite x from the parent scope?

// todo: infinite recursive calls?
/**
let Node = func => struct{
    left: Node();
    Right: Node();
};
*/

// todo: break makes parent loop evaluate to an expression?
// let index = for i in [10, 20, 30] => if i == 20 => break i;
// println(x); // 1
// but could be replaced by a function with a return:
// let index = (func => for i in [10, 20, 30] => if i == 20 => return i)();

// todo: use
// use import("std.txt");
// range(10); // no need to write "std." prefix

// todo: type checks
// let p = func (x: number, y: number) => x*y;
// p("qwe"); // error

// todo: return multiple values and multi-assign
// let x, y = func() => 42, 45;

// todo: slices?
// let s = "12345".(1:3); // "23"

// todo: fix: std.repeatcb(func=>std.repeat(".", 1000), 500): parsed as std."repeatcb"(func() return => tuple(std."repeat"(".", 1000), 500))

// todo: maybe resolve ints to ranges in for loops?
// e.g.: for i in len(a) => println(i);

// todo: cyclic dependencies in imports?
/**
e.g.:
let module = import("2022_13.txt");
module.solve()

<std.txt>
exports { // defines public interface
	version: "0.1"; // instant declaration
	sort;           // symbol is declared later
	range;
	Set;
	fill: func(x, cnt) {...};
};

let sort = func(...){...};
let Set = func => struct {...};



<other.txt>
let std = import("std.txt");
std.sort(array, less);
std.range(10);
*/

func New(r io.Reader, file string) *Lexer {
	l := &Lexer{
		r: bufio.NewReader(r),
		loc: &Location{
			File:   file,
			Column: 1,
			Line:   1,
		},
	}
	return l
}
func NewFromString(input string) *Lexer {
	return New(strings.NewReader(input), "(string input)")
}
func (l *Lexer) cur() rune {
	if len(l.buf) == 0 {
		return 0
	}
	return l.buf[0]
}
func (l *Lexer) reserve(n int) bool {
	for !l.eof && len(l.buf) < n {
		r, n, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.eof = true
				break
			} else {
				panic(err)
			}
		}
		if n == 0 {
			break
		}
		l.buf = append(l.buf, r)
	}
	return len(l.buf) >= n
}
func (l *Lexer) advance() bool {
	l.loc.Column += 1
	if l.cur() == '\n' {
		l.loc.Line += 1
		l.loc.Column = 1
	}
	if len(l.buf) >= 1 {
		l.buf = l.buf[1:]
	}
	l.reserve(1)
	return true
}
func (l *Lexer) advanceN(n int) bool {
	for i := 0; i < n; i++ {
		if !l.advance() {
			return false
		}
	}
	return true
}
func (l *Lexer) lookahead(needle []rune) bool {
	if !l.reserve(len(needle)) {
		return false
	}
	return string(l.buf[:len(needle)]) == string(needle)
}
func (l *Lexer) consume(needle []rune) bool {
	if !l.lookahead(needle) {
		return false
	}
	l.prevLoc = l.loc.Clone()
	l.advanceN(len(needle))
	return true
}
func (l *Lexer) consumeUntil(barrier []rune) (string, bool) {
	result := ""
	for l.cur() != 0 && !l.lookahead(barrier) {
		result += string(l.cur())
		l.advance()
	}

	return result, !l.consume(barrier)
}
func (l *Lexer) readIdentifier() string {
	result := ""
	for unicode.IsLetter(l.cur()) || unicode.IsNumber(l.cur()) || l.cur() == '_' {
		result += string(l.cur())
		l.advance()
	}
	return result
}
func (l *Lexer) readNumber() (string, error) {
	loc := l.loc.Clone()
	result := ""
	alphabet := map[rune]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
	}
	base := 10
	if l.consume([]rune("0b")) {
		alphabet = map[rune]bool{
			'0': true,
			'1': true,
		}
		base = 2
	} else if l.consume([]rune("0x")) {
		alphabet = map[rune]bool{
			'0': true,
			'1': true,
			'2': true,
			'3': true,
			'4': true,
			'5': true,
			'6': true,
			'7': true,
			'8': true,
			'9': true,
			'a': true,
			'b': true,
			'c': true,
			'd': true,
			'e': true,
			'f': true,
		}
		base = 16
	}

	for alphabet[l.cur()] {
		result += string(l.cur())
		l.advance()
	}
	if unicode.IsLetter(l.cur()) || unicode.IsDigit(l.cur()) {
		return "", fmt.Errorf("%s: %w", loc.String(), ErrInvalidNumberFormat) // number followed by an invalid string, e.g.: 123zzz
	}
	if result == "" {
		return "", fmt.Errorf("%s: %w", loc.String(), ErrInvalidNumberFormat) // prefix followed by an invalid string, e.g.: 0xzzz
	}
	i, err := strconv.ParseInt(result, base, 64)
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w", loc.String(), err.Error(), ErrInvalidNumberFormat)
	}
	return strconv.FormatInt(i, 10), nil
}
func (l *Lexer) readSingleRune() string {
	result := string(l.cur())
	l.advance()
	return result
}
func (l *Lexer) Next() (*Token, error) {
	l.reserve(1)
	for unicode.IsSpace(l.cur()) {
		l.advance()
	}
	if l.cur() == 0 {
		//return nil, ErrEof
		return &Token{Kind: TokenTypeEof, Location: l.loc.Clone()}, nil
	}

	if l.consume([]rune("//")) {
		loc := l.prevLoc
		comment, _ := l.consumeUntil([]rune("\n"))
		return &Token{Kind: TokenTypeComment, Literal: comment, Location: loc}, nil
	}
	if l.consume([]rune("/*")) {
		loc := l.prevLoc
		comment, eof := l.consumeUntil([]rune("*/"))
		if eof {
			return nil, fmt.Errorf("%s: %w", loc.String(), ErrUnterminatedCommentBlock)
		}
		return &Token{Kind: TokenTypeComment, Literal: comment, Location: loc}, nil
	}
	if l.consume([]rune("\"")) {
		loc := l.prevLoc
		str, eof := l.consumeUntil([]rune("\""))
		if eof {
			return nil, fmt.Errorf("%s: %w", loc.String(), ErrUnterminatedString)
		}
		return &Token{Kind: TokenTypeString, Literal: str, Location: loc}, nil
	}

	// try tokens from longest to shortest
	for _, token := range tokens {
		// todo: false detections possible, e.g. "returnq(1)" will be parsed as "return q(1)"
		// probably need to consider word boundary
		// also similar problem: "123qwe" will be parsed as [123, qwe]
		if l.consume([]rune(token.literal)) {
			return &Token{
				Kind:     token.kind,
				Literal:  token.literal,
				Location: l.prevLoc,
			}, nil
		}
	}
	if unicode.IsLetter(l.cur()) {
		loc := l.loc.Clone()
		identifier := l.readIdentifier()
		var kind TokenKind
		var ok bool
		kind, ok = keywords[identifier]
		if !ok {
			kind = TokenTypeIdentifier
		}
		return &Token{
			Kind:     kind,
			Literal:  identifier,
			Location: loc,
		}, nil
	}
	if unicode.IsDigit(l.cur()) {
		loc := l.loc.Clone()
		n, err := l.readNumber()
		return &Token{
			Kind:     TokenTypeNumber,
			Literal:  n,
			Location: loc,
		}, err
	}

	return nil, fmt.Errorf("%s [%d]: %w", string(l.cur()), l.cur(), ErrUnknownByte)
}
