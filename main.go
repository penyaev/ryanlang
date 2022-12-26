package main

import (
	"fmt"
	"os"
	"ryanlang/eval"
	"ryanlang/lexer"
	"ryanlang/object"
	"ryanlang/parser"
)

func run(fn string) (object.Object, *eval.Evaluator) {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	l := lexer.New(f, fn)
	p := parser.New(l)
	e := eval.New()

	mod := p.ReadModule(fn)

	evaled := e.Eval(mod)
	return evaled, e
}

func main() {
	if len(os.Args) != 2 {
		panic("filename expected as the single argument")
	}
	evaled, _ := run(os.Args[1])

	if evaled.Type() == object.ERROR {
		for i, e := range evaled.(*object.Error).Last(4) {
			if i > 0 {
				fmt.Print("\t")
			}
			fmt.Printf("%s: %s\n", e.Loc.String(), e.Msg)
		}
	}

	//fmt.Printf("(done in %s)", time.Now().Sub(now).String())
}
