package lexer

import (
	"errors"
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	tc := []struct {
		s       string
		tks     []Token
		wantErr bool
	}{
		{s: "aaa+bbb;", tks: []Token{{
			Kind:    TokenTypeIdentifier,
			Literal: "aaa",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 1,
			},
		}, {
			Kind:    TokenTypePlus,
			Literal: "+",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 4,
			},
		}, {
			Kind:    TokenTypeIdentifier,
			Literal: "bbb",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 5,
			},
		}, {
			Kind:    TokenTypeSemicolon,
			Literal: ";",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 8,
			},
		}}},
		{s: "let lets = returns;", tks: []Token{{
			Kind:    TokenTypeLet,
			Literal: "let",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 1,
			},
		}, {
			Kind:    TokenTypeIdentifier,
			Literal: "lets",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 5,
			},
		}, {
			Kind:    TokenTypeAssign,
			Literal: "=",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 10,
			},
		}, {
			Kind:    TokenTypeIdentifier,
			Literal: "returns",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 12,
			},
		}, {
			Kind:    TokenTypeSemicolon,
			Literal: ";",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 19,
			},
		}}},
		{s: "q1we qwe1 0xff 0b100 0 qwe let0 qwe_rty", tks: []Token{{
			Kind:    TokenTypeIdentifier,
			Literal: "q1we",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 1,
			},
		}, {
			Kind:    TokenTypeIdentifier,
			Literal: "qwe1",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 6,
			},
		}, {
			Kind:    TokenTypeNumber,
			Literal: "255",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 11,
			},
		}, {
			Kind:    TokenTypeNumber,
			Literal: "4",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 16,
			},
		}, {
			Kind:    TokenTypeNumber,
			Literal: "0",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 22,
			},
		}, {
			Kind:    TokenTypeIdentifier,
			Literal: "qwe",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 24,
			},
		}, {
			Kind:    TokenTypeIdentifier,
			Literal: "let0",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 28,
			},
		}, {
			Kind:    TokenTypeIdentifier,
			Literal: "qwe_rty",
			Location: &Location{
				File:   "(string input)",
				Line:   1,
				Column: 33,
			},
		}}},
	}

	for _, tt := range tc {
		t.Run(tt.s, func(t *testing.T) {
			l := NewFromString(tt.s)
			for _, tk := range tt.tks {
				rtk, err := l.Next()
				if !tt.wantErr && err != nil {
					t.Errorf("unexpected error %s", err)
				}
				if tt.wantErr && err == nil {
					t.Errorf("expected error")
				}

				if !reflect.DeepEqual(&tk, rtk) {
					t.Errorf("want=%v, got=%v", tk, rtk)
				}
			}
			if !tt.wantErr && !l.eof {
				t.Errorf("eof expected in lexer")
			}
		})
	}
}
func TestLexer_Errors(t *testing.T) {
	tc := []struct {
		s       string
		wantErr error
	}{
		{"/*", ErrUnterminatedCommentBlock},
		{"\"", ErrUnterminatedString},
		{"0qwe", ErrInvalidNumberFormat},
		{"0xabcp", ErrInvalidNumberFormat},
		{"0b11999", ErrInvalidNumberFormat},
		{"0b", ErrInvalidNumberFormat},
		{"0x", ErrInvalidNumberFormat},
		{"'123", ErrUnknownByte},
	}

	for _, tt := range tc {
		t.Run(tt.s, func(t *testing.T) {
			l := NewFromString(tt.s)
			tks, err := l.Next()
			if err == nil {
				t.Fatalf("error expected, got no error, tokens=%v", tks)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error expected: %s, got %s", tt.wantErr.Error(), err.Error())
			}
		})
	}
}
