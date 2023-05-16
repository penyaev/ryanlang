package vm

import (
	"fmt"
	"ryanlang/compiler"
	"ryanlang/compiler/instruction"
	"ryanlang/object"
	"strings"
)

type breakpointEvent int

const (
	breakpointStep breakpointEvent = iota
	breakpointLeaveFrame
	breakpointDebuggerCall
	breakpointRuntimeError
)

type breakpointLocation struct {
	code *object.Code
	addr int
}
type breakpoints struct {
	events    uint64
	stop      bool
	locations map[*object.Code]map[int]bool
}

func (b *breakpoints) toggleLocation(code *object.Code, cp int) {
	if b.locations == nil {
		b.locations = map[*object.Code]map[int]bool{}
	}
	if b.locations[code] == nil {
		b.locations[code] = map[int]bool{}
	}
	b.locations[code][cp] = !b.locations[code][cp]
}
func (b *breakpoints) set(events ...breakpointEvent) {
	for _, event := range events {
		b.events |= 1 << event
	}
}
func (b *breakpoints) at(code *object.Code, cp int) {
	if _, ok := b.locations[code]; ok && b.locations[code][cp] {
		b.stop = true
	}
}
func (b *breakpoints) trigger(event breakpointEvent) {
	if b.events&(1<<event) != 0 {
		b.stop = true
	}
}
func (b *breakpoints) halt() bool {
	return b.stop
}
func (b *breakpoints) reset() {
	b.stop = false
	b.events = 0
}

func (v *VM) dumpCode(c *object.Code, indent int, dd *compiler.DebugData) {
	code := c.Code
	nextAnnotation := ""
	var cf *Frame
	var cfi int
	for fi := 0; fi <= v.fp; fi++ {
		f := v.frames[fi]
		if c == f.cl.Code {
			cf = f
			cfi = fi
		}
	}

	indentFormat := fmt.Sprintf("%% %ds ", indent*16-1)
	if cf != nil {
		fmt.Printf(indentFormat+"frame %d (%p, cl %p, code %p), rs=%s, %d foreigns, cp=%d, cpe=%d\n", "---", cfi-v.fp, cf, cf.cl, cf.cl.Code, cf.cl.Code.ReturnScope.String(), len(cf.cl.Foreigns), cf.cp, cf.cpe)
		for oi, obj := range cf.cl.Foreigns {
			fmt.Printf(indentFormat+"%04d\t%s\t%s\n", "F", oi, (*obj).Type().String() /*(*obj).String()*/, "<redacted>")
		}
		//if len(cf.cl.Foreigns) > 0 {
		fmt.Printf(indentFormat+"%s\n", "", "--------------------------------")
		//}
		fmt.Println()
	}

	i := 0
	//indentString := strings.Repeat("\t", indent)

	for i < len(code) {
		inst, n, err := instruction.Read(code[i:])
		if err != nil {
			panic(fmt.Errorf("reading instruction: %w", err))
		}
		/*if inst.Op() == instruction.OpAnnotation {
			obj, ok := v.objects.Get(inst.(instruction.Annotation).Index)
			if ok {
				nextAnnotation = obj.(*object.String).Value
			} else {
				nextAnnotation = "! annotation const index out of range"
			}
		} else*/{
			var comment string
			switch inst := inst.(type) {
			case instruction.Annotation:
				obj, _ := v.objects.Get(inst.Index)
				comment = fmt.Sprintf("\t%s", obj.(*object.String).Value)
			case instruction.PushConstant:
				var constValue string

				obj, ok := v.objects.Get(inst.Index)
				if ok {
					constValue = obj.String()
				} else {
					constValue = "! const index out of range"
				}

				comment = fmt.Sprintf("\t; %s", constValue)
			case instruction.StoreLocal:
				comment = fmt.Sprintf("\t; \"%s\"", dd.LocalName(int(inst.Index)))
			case instruction.PushLocal:
				comment = fmt.Sprintf("\t; \"%s\"", dd.LocalName(int(inst.Index)))
			case instruction.Jnt:
				var na int
				switch inst.AddrType {
				case compiler.RelativeAddress:
					na = i + n + int(inst.Addr)
				case compiler.RelativeAddressBackwards:
					na = i + n - int(inst.Addr)
				case compiler.AbsoluteAddress:
					na = int(inst.Addr)
				case compiler.LabelAddress:
					if cf != nil {
						na = cf.labels[instruction.LabelKind(inst.Addr)]
					}
				}
				comment = fmt.Sprintf("\t; --> %04d", na)
			case instruction.Jmp:
				var na int
				switch inst.AddrType {
				case compiler.RelativeAddress:
					na = i + n + int(inst.Addr)
				case compiler.RelativeAddressBackwards:
					na = i + n - int(inst.Addr)
				case compiler.AbsoluteAddress:
					na = int(inst.Addr)
				case compiler.LabelAddress:
					if cf != nil {
						na = cf.labels[instruction.LabelKind(inst.Addr)]
					}
				}
				comment = fmt.Sprintf("\t; --> %04d", na)
			}
			if nextAnnotation != "" {
				comment += " <" + nextAnnotation + ">"
				nextAnnotation = ""
			}

			var marks []string
			if cf != nil && cf.cpe == i {
				marks = append(marks, "cp"+strings.Repeat("@", v.fp-cfi))

				entry := dd.SearchSource(i)
				if entry != nil {
					comment = fmt.Sprintf("\t; %s at %s", entry.Comment, entry.Loc.String())
				}
			}

			left := strings.Join(marks, ", ")
			if left != "" {
				left = left + " >"
			}

			fmt.Printf(indentFormat+"%04d\t%s%s\n", left, i, inst.String(), comment)
		}
		i += n
	}
}
func (v *VM) Dump() {
	fmt.Println("code:")
	v.dumpCode(v.frames[0].cl.Code, 1, nil)

	fmt.Println()
	fmt.Println("objects:")
	for i := uint16(0); i < v.objects.Len(); i++ {
		obj, _ := v.objects.Get(i)
		if obj.Type() == object.CODE {
			fmt.Println()
		}
		fmt.Printf("\t%04d\t%s\t%s\n", i, obj.Type().String(), obj.String())
		if obj.Type() == object.CODE {
			v.dumpCode(obj.(*object.Code), 2, v.debugData[int(i)])
		}
	}
	/*
		fmt.Println()
		fmt.Println("frames:")
		for fi := 0; fi <= v.fp; fi++ {
			f := v.frames[fi]

			//fmt.Printf("frame %d\n", fi-v.fp)
			fmt.Println()
			v.dumpCode(f.cl.Code, 1)
		}*/
}
func (v *VM) dumpStack() {
	fmt.Println("stack:")
	for si := v.sp; si >= 0; si-- {
		value := *v.stack[si]

		var marks []string
		if si == v.sp {
			if v.frame != nil && v.sp == v.frame.bsp+v.frame.cl.Code.Locals {
				fmt.Printf("% 31s %s\n", "sp  ", "------------------------------")
			} else {
				marks = append(marks, "sp")
			}
		}
		if v.frame != nil && si == v.frame.bsp+1 {
			marks = append(marks, "bsp")
		}
		if v.frame != nil && si > v.frame.bsp && si <= v.frame.bsp+v.frame.cl.Code.Locals {
			marks = append(marks, fmt.Sprintf("loc #%d", si-v.frame.bsp-1))
		}
		left := strings.Join(marks, ", ")
		if left != "" {
			left = left + " >"
		}

		fmt.Printf("% 31s %04d\t%s\t%s\n", left, si, value.Type().String() /*value.String()*/, "")
	}
}
