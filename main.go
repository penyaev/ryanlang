package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"ryanlang/compiler"
	"ryanlang/eval"
	"ryanlang/lexer"
	"ryanlang/object"
	"ryanlang/parser"
	"ryanlang/vm"
	"time"
)

type Engine int

const (
	Eval Engine = iota
	VM
)

func run(fn string, engine Engine) (object.Object, *eval.Evaluator) {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	l := lexer.New(f, fn)
	p := parser.New(l)

	mod := p.ReadModule(fn)
	var evaled object.Object

	if engine == VM {
		c := compiler.NewCompiler()
		compiledModule, err := c.CompileRunModule(mod)
		if err != nil {
			panic(fmt.Errorf("compilation error: %w", err))
		}

		v := vm.New(compiledModule)
		v.EnableWebDebugger()
		if fnBytes, err := os.ReadFile(fn); err == nil {
			v.AddSourceFile(fn, string(fnBytes))
		} else {
			panic(err)
		}

		now := time.Now()
		evaled = v.Run()
		fmt.Println(time.Now().Sub(now))
	} else if engine == Eval {
		e := eval.New()
		now := time.Now()
		evaled = e.Eval(mod)
		fmt.Println(time.Now().Sub(now))
	} else {
		panic("unknown engine")
	}

	return evaled, nil
}

func main() {
	if len(os.Args) != 2 {
		panic("filename expected as the single argument")
	}

	f, err := os.Create("cpuprofile")
	if err != nil {
		panic(err)
	}
	if err = pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}

	//hf, err := os.Create("heapprofile")
	//if err != nil {
	//	panic(err)
	//}

	//debug.SetGCPercent(-1)
	evaled, _ := run(os.Args[1], VM)
	//runtime.GC()
	//if err := pprof.WriteHeapProfile(hf); err != nil {
	//	panic(err)
	//}
	pprof.StopCPUProfile()

	if evaled.Type() == object.ERROR {
		for i, e := range evaled.(*object.Error).Last(4) {
			if i > 0 {
				fmt.Print("\t")
			}
			fmt.Printf("%s: %s\n", e.Loc.String(), e.Msg)
		}
	} else {
		fmt.Println(evaled.String())
	}

	//fmt.Printf("(done in %s)", time.Now().Sub(now).String())
}
