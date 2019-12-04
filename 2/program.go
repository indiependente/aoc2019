package main

import (
	"fmt"
	"strings"
)

type OpCode = int
type Instruction []OpCode

func (i Instruction) String() string {
	if i[0] == HaltOp {
		return fmt.Sprintf("%d", i[0])
	}
	return fmt.Sprintf("%d,%d,%d,%d", i[0], i[1], i[2], i[3])
}

type InstructionSet map[OpCode]func(int, int) (int, error)

type Program []OpCode

func (p Program) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", p[0]))
	for i := 1; i < len(p); i++ {
		sb.WriteString(",")
		sb.WriteString(fmt.Sprintf("%d", p[i]))
	}
	return sb.String()
}
func (p *Program) AddInstruction(i Instruction) {
	*p = append(*p, i...)
}
