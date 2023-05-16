package compiler

import (
	"ryanlang/compiler/sourcemap"
	"ryanlang/object"
)

func iferr(args ...interface{}) error {
	for _, arg := range args {
		e, ok := arg.(error)
		if ok {
			return e
		}
	}
	return nil
}

type scope struct {
	code *code
}

func newScope() *scope {
	return &scope{
		code: &code{
			sm: sourcemap.New(),
		},
	}
}

type Module struct {
	EntryPoint uint16
	Objects    *object.Storage
	DebugData  map[int]*DebugData
}
type code struct {
	b  []byte
	sm *sourcemap.SourceMap
}

func (c *code) Len() int {
	return len(c.b)
}
