package vm

import (
	"ryanlang/compiler/instruction"
	"ryanlang/object"
)

type Frame struct {
	cl     *object.Closure
	cp     int
	cpe    int
	bsp    int
	labels map[instruction.LabelKind]int
}
