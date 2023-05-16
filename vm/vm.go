package vm

import (
	"bufio"
	"fmt"
	"os"
	"ryanlang/compiler"
	"ryanlang/compiler/instruction"
	"ryanlang/funcs"
	"ryanlang/lexer"
	"ryanlang/object"
	"ryanlang/parser"
	"strings"
)

var ErrDebuggerHalt = fmt.Errorf("debugger halt")

type state int

const (
	statePaused state = iota
	stateRunning
)

type VM struct {
	stack       []*object.Object
	sp          int // points at the top of the stack, -1 means it points below the stack's beginning, i.e. stack is empty
	frames      []*Frame
	frame       *Frame
	fp          int
	objects     *object.Storage
	debugData   map[int]*compiler.DebugData
	sourceFiles map[string]string
	state       state
	wd          *webDebugger
	bp          *breakpoints
}

func New(compiledModule *compiler.Module) *VM {
	v := &VM{
		sp:          -1,
		fp:          -1,
		objects:     compiledModule.Objects,
		debugData:   compiledModule.DebugData,
		sourceFiles: map[string]string{},
		state:       statePaused,
		bp:          &breakpoints{},
	}
	entrypoint, ok := compiledModule.Objects.Get(compiledModule.EntryPoint)
	if !ok {
		panic("cannot find entrypoint object in module")
	}
	v.enterFrame(&object.Closure{
		Code: entrypoint.(*object.Code),
	})

	return v
}

func (v *VM) enterFrame(cl *object.Closure) {
	v.fp++

	if v.fp >= len(v.frames) {
		nf := make([]*Frame, len(v.frames)*2+1)
		copy(nf, v.frames)
		v.frames = nf
	}
	if v.frames[v.fp] == nil {
		v.frames[v.fp] = &Frame{}
	}
	// reusing existing frame objects
	v.frames[v.fp].cl = cl
	v.frames[v.fp].bsp = v.sp - cl.Code.Arguments
	v.frames[v.fp].cp = 0
	v.frames[v.fp].cpe = 0
	if v.frames[v.fp].labels != nil && len(v.frames[v.fp].labels) != 0 {
		v.frames[v.fp].labels = nil
	}
	v.frame = v.frames[v.fp]

	for i := 0; i < cl.Code.Locals-cl.Code.Arguments; i++ { // reserve space on the stack for local vars
		v.pushNull()
	}
}
func (v *VM) leaveFrame() *Frame {
	// todo: check if there's no frames
	// todo: check if there's more than 1 value on the stack, it might be a leak
	if v.fp == -1 {
		panic("frame stack empty, cannot leave frame")
	}
	f := v.frame
	v.fp--

	if v.fp >= 0 {
		v.frame = v.frames[v.fp]
	} else {
		v.frame = nil
	}

	v.sp = f.bsp // reset stack to where it was before calling the frame

	v.bp.trigger(breakpointLeaveFrame)
	return f
}
func (v *VM) returnFromFrame(scope object.CodeReturnScope) error {
	if v.frame == nil {
		return fmt.Errorf("frame stack is empty")
	}

	if scope == object.CodeReturnScopeContinue {
		if v.frame.cl.Code.ReturnScope != object.CodeReturnScopeLoop {
			return fmt.Errorf("cannot continue because not inside a loop")
		}
		v.frame.cp = 0
		v.frame.cpe = 0
	}

	for {
		f := v.leaveFrame()
		if f.cl.Code.ReturnScope == scope {
			break
		}
		if f.cl.Code.ReturnScope < scope {
			return fmt.Errorf("cannot return because there was no appropriate scope to return from")
		}
	}
	return nil
}

func (v *VM) push(obj *object.Object) {
	v.sp++

	if v.sp > 10000 {
		panic("stack overflow")
	}

	if v.sp >= len(v.stack) {
		ns := make([]*object.Object, 2*len(v.stack)+1)
		copy(ns, v.stack)
		v.stack = ns
	}

	v.stack[v.sp] = obj
}
func (v *VM) pushNull() {
	var n object.Object = &object.StaticNull
	v.push(&n)
}

func (v *VM) expectPop(typ object.Type) (*object.Object, error) {
	value := v.pop()
	if value == nil {
		return nil, fmt.Errorf("unexpected nil on the stack")
	}
	if (*value).Type() != typ {
		return nil, fmt.Errorf("expected: %s, got: %s", typ.String(), (*value).Type().String())
	}
	return value, nil
}
func (v *VM) pop() *object.Object {
	if v.sp == v.frame.bsp+v.frame.cl.Code.Locals {
		v.Dump()
		v.dumpStack()
		panic(fmt.Errorf("stack underflow sp=%d", v.sp))
	}

	v.sp--
	return v.stack[v.sp+1]
}
func (v *VM) top() *object.Object {
	if v.sp == -1 {
		panic("getting top of empty stack")
		//var ret object.Object = &object.StaticNull
		//return &ret // todo: check if this is correct? returning null when the stack is empty
	}
	if v.sp < 0 || v.sp >= len(v.stack) {
		panic("broken stack? or just empty")
	}

	return v.stack[v.sp]
}
func (v *VM) setState(state state) {
	v.state = state
}
func (v *VM) next() (bool, error) {
	v.setState(stateRunning)
	defer v.setState(statePaused)

	args := make([]interface{}, 2)
	var c []byte
	var op instruction.Op
	var n int
	binaryOps := map[instruction.Op]func(left object.Object, right object.Object) object.Object{
		instruction.OpAdd:         funcs.Plus,
		instruction.OpSub:         funcs.Minus,
		instruction.OpMult:        funcs.Mult,
		instruction.OpDiv:         funcs.Div,
		instruction.OpMod:         funcs.Mod,
		instruction.OpGt:          funcs.Gt,
		instruction.OpGte:         funcs.Gte,
		instruction.OpLt:          funcs.Lt,
		instruction.OpLte:         funcs.Lte,
		instruction.OpEqTest:      funcs.EqTest,
		instruction.OpFieldAccess: funcs.FieldAccess,
		instruction.OpLogicalOr:   funcs.LogicalOr,
	}
	for {
		if v.frame == nil {
			// todo: no frames left, no code left, nothing to do?
			return false, nil
		}
		c = v.frame.cl.Code.Code
		if v.frame.cp >= len(c) {
			// todo: end of code reached, possibly return missing?
			//fmt.Printf("leaving frame, top: %s\n", (*v.top()).String())
			t := v.top()
			v.leaveFrame()
			v.push(t) // frame will resolve to whatever what on top of the stack when we were leaving
			//v.pushNull()
			//return true, nil
			continue
		}
		//inst, n, err := instruction.Read(c[v.frame().cp:])

		op, n = instruction.ReadFast(c, v.frame.cp, args)

		//if err != nil {
		//	panic(fmt.Errorf("reading instruction: %w", err))
		//}
		v.frame.cpe = v.frame.cp
		v.frame.cp += n

		switch op {
		default:
			handler, ok := binaryOps[op]
			if ok {
				r := handler(*v.pop(), *v.pop())
				if object.IsError(r) {
					return false, fmt.Errorf("operator %s: %s", op.String(), r.(*object.Error).String())
				}
				v.push(&r)
			} else {
				panic("dont know how to run op: " + op.String())
			}
		case instruction.OpPushConstant:
			obj, ok := v.objects.GetRef(args[0].(uint16))
			if !ok {
				return false, fmt.Errorf("constant index out of range")
			}
			v.push(obj)
		case instruction.OpJnt:
			addrType := args[0].(uint8)
			addr := args[1].(uint16)
			top, err := v.expectPop(object.BOOLEAN)
			if err != nil {
				return false, fmt.Errorf("operator %s: %w", op.String(), err)
			}
			if (*top).(*object.Boolean).Value == false {
				if addrType == compiler.RelativeAddress {
					v.frame.cp += int(addr)
				} else if addrType == compiler.AbsoluteAddress {
					v.frame.cp = int(addr)
				} else if addrType == compiler.RelativeAddressBackwards {
					v.frame.cp -= int(addr)
				} else if addrType == compiler.LabelAddress {
					dest, ok := v.frame.labels[instruction.LabelKind(addr)]
					if !ok {
						return false, fmt.Errorf("label is not set")
					}
					v.frame.cp = dest
				} else {
					panic("unexpected address type")
				}
			}
		case instruction.OpJmp:
			addrType := args[0].(uint8)
			addr := args[1].(uint16)
			// todo: move read address into separate func
			// todo: move handle jump to an address into separate func
			if addrType == compiler.RelativeAddress {
				v.frame.cp += int(addr)
			} else if addrType == compiler.AbsoluteAddress {
				v.frame.cp = int(addr)
			} else if addrType == compiler.RelativeAddressBackwards {
				v.frame.cp -= int(addr)
			} else if addrType == compiler.LabelAddress {
				dest, ok := v.frame.labels[instruction.LabelKind(addr)]
				if !ok {
					return false, fmt.Errorf("label is not set")
				}
				v.frame.cp = dest
			} else {
				panic("unexpected address type")
			}
		case instruction.OpAnnotation:
			// ignore
		case instruction.OpCall:
			obj, err := v.expectPop(object.CLOSURE)
			if err != nil {
				return false, fmt.Errorf("call: %w", err)
			}
			callee := (*obj).(*object.Closure)
			if callee.BuiltinFunctionName != "" {
				builtin, ok := funcs.BuiltinFunctions[callee.BuiltinFunctionName]
				if !ok {
					return false, fmt.Errorf("unknown builting function name: %s", callee.BuiltinFunctionName)
				}
				if len(builtin.Arguments) != int(args[0].(uint8)) {
					return false, fmt.Errorf("built-in %s: expected %d arguments, got %d", callee.BuiltinFunctionName, len(builtin.Arguments), args[0].(uint8))
				}
				argsmap := map[string]object.Object{}
				for i := len(builtin.Arguments) - 1; i >= 0; i-- {
					argsmap[builtin.Arguments[i]] = *v.pop()
				}

				var ret object.Object
				ret = builtin.Body(argsmap)
				switch ret := ret.(type) {
				case *object.ReturnObject:
					v.push(&ret.Obj)
				case *object.Error:
					return false, fmt.Errorf("built-in %s: %s", callee.BuiltinFunctionName, ret.String())
				}
				if callee.BuiltinFunctionName == "debugger" {
					v.bp.trigger(breakpointDebuggerCall)
				}
			} else {
				if callee.Code.Arguments != int(args[0].(uint8)) {
					return false, fmt.Errorf("expected %d arguments, got %d", callee.Code.Arguments, args[0].(uint8))
				}
				v.enterFrame(callee)
				// todo: push onto the stack whatever the called function returned
			}
		case instruction.OpArray:
			itemsc := int(args[0].(uint16))
			array := &object.Array{
				make([]object.Object, itemsc),
			}
			for i := 0; i < itemsc; i++ {
				array.Items[i] = *v.pop()
			}
			var obj object.Object = array
			v.push(&obj)
		case instruction.OpStruct:
			itemsc := int(args[0].(uint8))
			str := &object.Struct{
				Fields: make(map[string]object.Object, itemsc),
			}
			for i := 0; i < itemsc; i++ {
				k, err := v.expectPop(object.STRING)
				if err != nil {
					return false, err
				}
				v := v.pop()
				str.Fields[(*k).(*object.String).Value] = *v
			}
			var obj object.Object = str
			v.push(&obj)
		case instruction.OpMap:
			itemsc := int(args[0].(uint8))
			m := &object.Map{
				Fields: make(map[string]object.MapItem, itemsc),
			}
			for i := 0; i < itemsc; i++ {
				k := v.pop()
				if !object.IsHashable(*k) {
					return false, fmt.Errorf("hashable map key expected, got: %s", (*k).Type().String())
				}
				m.Fields[(*k).(object.Hashable).Hash()] = object.MapItem{
					Key:   *k,
					Value: *v.pop(),
				}
			}
			var obj object.Object = m
			v.push(&obj)
		case instruction.OpTuple:
			itemsc := int(args[0].(uint8))
			t := &object.Tuple{
				Values: make([]object.Object, itemsc),
			}
			for i := 0; i < itemsc; i++ {
				t.Values[i] = *v.pop()
			}
			var obj object.Object = t
			v.push(&obj)
		case instruction.OpUntuple:
			itemsc := int(args[0].(uint8))
			t, err := v.expectPop(object.TUPLE)
			if err != nil {
				return false, fmt.Errorf("untuple: %w", err)
			}
			if len((*t).(*object.Tuple).Values) != itemsc {
				return false, fmt.Errorf("untuple: %d items expected, got %d", itemsc, len((*t).(*object.Tuple).Values))
			}
			for i := len((*t).(*object.Tuple).Values) - 1; i >= 0; i-- {
				value := (*t).(*object.Tuple).Values[i]
				v.push(&value)
			}
		case instruction.OpStoreLocal:
			*v.stack[v.frame.bsp+1+int(args[0].(uint16))] = *v.top() // don't pop because this op should resolve to the assigned value
		case instruction.OpStoreForeign:
			*v.frame.cl.Foreigns[int(args[0].(uint16))] = *v.top() // don't pop because this op should resolve to the assigned value
		case instruction.OpPushLocalRef:
			// when passed like this, whoever uses this variable, they will have full control over it,
			// it will be able to rewrite its value
			// essentially it's like if this variable was in a closure:
			// let f = func(x) { x++; }; let a = 1; f(a); // <-- a will now be 2.
			if int(args[0].(uint16)) > v.frame.cl.Code.Locals {
				return false, fmt.Errorf("trying to get local %d in a closure with %d locals", int(args[0].(uint16)), v.frame.cl.Code.Locals)
			}
			v.push(v.stack[v.frame.bsp+1+int(args[0].(uint16))])

			// if needed to pass a copy, use this instead (create a new pointer to the same value):
			//val := *v.stack[v.frame.bsp+1+int(args[0].(uint16))]
			//v.push(&val)
		case instruction.OpCopy:
			// new pointer now points to the same value (see comments above)
			val := *v.top()
			v.stack[v.sp] = &val
		case instruction.OpDup:
			v.push(v.top())
		case instruction.OpPushForeign:
			v.push(v.frame.cl.Foreigns[int(args[0].(uint16))])
		case instruction.OpPop:
			v.pop()
		case instruction.OpReturn:
			scope := object.CodeReturnScope(args[0].(uint8))
			ret := v.pop()
			err := v.returnFromFrame(scope)
			if err != nil {
				return false, fmt.Errorf("return: %w", err)
			}
			v.push(ret)
		case instruction.OpClosure:
			obj, err := v.expectPop(object.CODE)
			if err != nil {
				return false, fmt.Errorf("closure: %w", err)
			}
			objTyped := (*obj).(*object.Code)
			cl := &object.Closure{
				Code:     objTyped,
				Foreigns: make([]*object.Object, objTyped.Foreigns),
			}
			for i := 0; i < objTyped.Foreigns; i++ {
				cl.Foreigns[i] = v.pop()
			}
			var cli object.Object = cl
			v.push(&cli)
		case instruction.OpFieldAssign:
			left := *v.pop()
			right := *v.pop()
			value := *v.pop()
			ret := funcs.FieldAssign(left, right, value)
			if object.IsError(ret) {
				return false, fmt.Errorf("field assign: %s", ret.String())
			}
			v.push(&ret)
		case instruction.OpImport:
			modulePath, err := v.expectPop(object.STRING)
			if err != nil {
				return false, fmt.Errorf("import: %w", err)
			}

			fn := (*modulePath).(*object.String).Value
			f, err := os.Open(fn)
			if err != nil {
				return false, fmt.Errorf("import: cannot open file: %w", err)
			}
			l := lexer.New(f, fn)
			p := parser.New(l)

			mod := p.ReadModule(fn)
			err = f.Close()
			if err != nil {
				return false, fmt.Errorf("import: cannot close file: %w", err)
			}

			compiledModule, err := compiler.NewCompilerWithStorage(v.objects).CompileModule(mod)
			if err != nil {
				return false, fmt.Errorf("import %s: module compilation error: %w", fn, err)
			}

			for id, sm := range compiledModule.DebugData {
				v.debugData[id] = sm
			}
			if fnBytes, err := os.ReadFile(fn); err == nil {
				v.AddSourceFile(fn, string(fnBytes))
			} else {
				return false, fmt.Errorf("import: cannot read file: %w", err)
			}

			module, ok := v.objects.Get(compiledModule.EntryPoint)
			if !ok {
				return false, fmt.Errorf("import: cannot find module object")
			}
			v.push(&module)
		case instruction.OpLabel:
			kind := instruction.LabelKind(args[0].(uint8))
			if v.frame.labels == nil {
				v.frame.labels = map[instruction.LabelKind]int{}
			}
			v.frame.labels[kind] = v.frame.cp
		}
		if v.frame != nil {
			v.frame.cpe = v.frame.cp
			v.bp.at(v.frame.cl.Code, v.frame.cp)
		}

		v.bp.trigger(breakpointStep)
		if v.bp.halt() {
			break
		}
	}
	return true, nil
}
func (v *VM) AddSourceFile(fn string, contents string) {
	v.sourceFiles[fn] = contents
}
func (v *VM) EnableWebDebugger() {
	v.wd = newWebDebugger(v)
	go v.wd.start()
}
func (v *VM) Run() object.Object {
	v.bp.set(breakpointDebuggerCall)
	v.bp.set(breakpointRuntimeError)
	for {
		more, err := v.next()
		stop := v.bp.halt() || err != nil || !more
		if stop {
			//if (more && step) || err == ErrDebuggerHalt {
			if err != nil {
				v.bp.trigger(breakpointRuntimeError)
				fmt.Printf("runtime error: %s\n", err.Error())
			}
			if v.bp.halt() {
				v.bp.reset()

				if v.wd != nil {
					v.wd.Notify(vmStateChanged, nil)
					if err != nil {
						v.wd.Notify(vmHaltMsg, err)
					}

					fmt.Printf("waiting for web debugger commands...\n")
					commands := v.wd.Listen()
					cont := false
					for !cont {
						cmd := <-commands
						cont = true
						switch cmd.typ {
						case webCommandNext:
							v.bp.set(breakpointStep)
						case webCommandRun:
							v.bp.set(breakpointDebuggerCall)
						case webCommandStepOut:
							v.bp.set(breakpointLeaveFrame)
						case webCommandToggleBreakpoint:
							addr := cmd.payload.(webCommandPayloadToggleBreakpoint).Address
							objID := cmd.payload.(webCommandPayloadToggleBreakpoint).ObjectID
							if obj, ok := v.objects.Get(uint16(objID)); ok && obj.Type() == object.CODE {
								v.bp.toggleLocation(obj.(*object.Code), addr)
								v.wd.Notify(vmBreakpointsUpdated, nil)
							}
							cont = false
						default:
							panic("unknown web debugger command")
						}
					}
					v.bp.set(breakpointRuntimeError)
					v.wd.StopListening()
				} else {
					fmt.Println()
					if err != nil {
						fmt.Println(err)
					}

					// todo: make cleaner
					// debugger
					v.Dump()
					v.dumpStack()

					cont := false
					def := "next"
					for !cont {
						reader := bufio.NewReader(os.Stdin)
						fmt.Printf("debugger [%s]> ", def)
						cmd, err := reader.ReadString('\n')
						if err != nil {
							panic(err)
						}

						cont = true
						cmd = strings.TrimSpace(cmd)
						if cmd == "" {
							cmd = def
						}
						switch cmd {
						case "next":
							v.bp.set(breakpointStep)
						case "run":
							v.bp.set(breakpointDebuggerCall)
						case "stepout":
							v.bp.set(breakpointLeaveFrame)
						default:
							fmt.Println("unknown command")
							cont = false
						}
					}
				}

			} else if err != nil {
				panic(err)
			} else if !more {
				break
			} else {
				panic("now what?")
			}
		}

		if !more {
			break
		}
	}

	return *v.top()
}
