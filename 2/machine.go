package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	AddOp                   OpCode       = 1
	MulOp                   OpCode       = 2
	HaltOp                  OpCode       = 99
	HaltSym                 string       = `99`
	ErrMalformedInstruction MachineError = `malformed instruction`
	ErrUnknownOpCode        MachineError = `unknown operation code`
	ErrBadMemoryAccess      MachineError = `unexpected access to memory location`
	EOX                     MachineError = `end of execution`
)

type MachineError string

func (e MachineError) Error() string {
	return string(e)
}

type OpCode int
type Instruction struct {
	Op    OpCode
	Term1 int
	Term2 int
	Store int
}

func (i Instruction) String() string {
	if i.Op == HaltOp {
		return fmt.Sprintf("%d", i.Op)
	}
	return fmt.Sprintf("%d,%d,%d,%d", i.Op, i.Term1, i.Term2, i.Store)
}

type InstructionSet map[OpCode]func(int, int) (int, error)

type Program []Instruction

func (p Program) String() string {
	var sb strings.Builder
	sb.WriteString(p[0].String())
	for i := 1; i < len(p); i++ {
		sb.WriteString(",")
		sb.WriteString(p[i].String())
	}
	return sb.String()
}
func (p *Program) AddInstruction(i Instruction) {
	*p = append(*p, i)
}

type IntcodeMachine struct {
	Program Program
	ISet    InstructionSet
}

func parseLine(symbols []string) (Instruction, int, error) {
	op, err := strconv.Atoi(symbols[0])
	if err != nil {
		return Instruction{}, 0, fmt.Errorf("could not parse instruction '%s' : %w", symbols[0], ErrMalformedInstruction)
	}
	if len(symbols) == 1 {
		if symbols[0] == HaltSym {
			return Instruction{
				Op:    OpCode(op),
				Term1: 0,
				Term2: 0,
				Store: 0,
			}, 1, nil
		} else {
			return Instruction{}, 0, fmt.Errorf("could not parse instruction '%s' : %w", symbols[0], ErrMalformedInstruction)
		}

	}
	term1, err := strconv.Atoi(symbols[1])
	if err != nil {
		return Instruction{}, 0, fmt.Errorf("could not parse instruction '%s' : %w", symbols[1], ErrMalformedInstruction)
	}
	term2, err := strconv.Atoi(symbols[2])
	if err != nil {
		return Instruction{}, 0, fmt.Errorf("could not parse instruction '%s' : %w", symbols[2], ErrMalformedInstruction)
	}
	store, err := strconv.Atoi(symbols[3])
	if err != nil {
		return Instruction{}, 0, fmt.Errorf("could not parse instruction '%s' : %w", symbols[3], ErrMalformedInstruction)
	}
	return Instruction{
		Op:    OpCode(op),
		Term1: term1,
		Term2: term2,
		Store: store,
	}, 4, nil
}

func NewIntcodeMachine(r io.Reader, sep string) (*IntcodeMachine, error) {
	program := Program{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate machine: %w", err)
	}
	symbols := strings.Split(string(data), sep)
	for i := 0; i < len(symbols); {
		var line []string
		if symbols[i] == HaltSym {
			line = []string{symbols[i]}
		} else {
			line = append(line, symbols[i], symbols[i+1], symbols[i+2], symbols[i+3])
		}
		instr, size, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("could not instantiate machine: %w", err)
		}
		program.AddInstruction(instr)
		i += size
	}
	return &IntcodeMachine{
		Program: program,
		ISet:    make(InstructionSet),
	}, nil
}
func NewIntcodeMachineWithStdISet(r io.Reader, sep string) (*IntcodeMachine, error) {
	m, err := NewIntcodeMachine(r, sep)
	if err != nil {
		return nil, err
	}
	err = m.AddInstruction(AddOp, func(a, b int) (int, error) {
		return a + b, nil
	})
	if err != nil {
		return nil, err
	}
	err = m.AddInstruction(MulOp, func(a, b int) (int, error) {
		return a * b, nil
	})
	err = m.AddInstruction(HaltOp, func(a, b int) (int, error) {
		return 0, EOX
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (icm *IntcodeMachine) fetch(addr int) (*int, error) {
	switch addr % 4 {
	case 0:
		return (*int)(&(icm.Program[addr/4].Op)), nil
	case 1:
		return &(icm.Program[addr/4].Term1), nil
	case 2:
		return &(icm.Program[addr/4].Term2), nil
	case 3:
		return &(icm.Program[addr/4].Store), nil
	default:
		return nil, fmt.Errorf("could not access location [%d,%d]: %w", addr/4, addr%4, ErrBadMemoryAccess)
	}
}

func (icm *IntcodeMachine) AddInstruction(code OpCode, fn func(int, int) (int, error)) error {
	icm.ISet[code] = fn
	return nil
}

func (icm *IntcodeMachine) executeInstruction(i Instruction) error {
	fn, ok := icm.ISet[i.Op]
	if !ok {
		return fmt.Errorf("could not execute instruction %s: %w", i, ErrUnknownOpCode)
	}
	a, err := icm.fetch(i.Term1)
	if err != nil {
		return fmt.Errorf("could not fetch address %d: %w", i.Term1, err)
	}
	b, err := icm.fetch(i.Term2)
	if err != nil {
		return fmt.Errorf("could not fetch address %d: %w", i.Term2, err)
	}
	v, err := fn(*a, *b)
	if err != nil {
		if err == EOX {
			return EOX
		}
		return fmt.Errorf("could not execute instruction %s", i)
	}
	out, err := icm.fetch(i.Store)
	if err != nil {
		return fmt.Errorf("could not fetch address %d: %w", i.Store, err)
	}
	*out = v
	return nil
}

func (icm *IntcodeMachine) Execute() error {
	for _, instr := range icm.Program {
		err := icm.executeInstruction(instr)
		if err != nil {
			if err == EOX {
				return EOX
			}
			return fmt.Errorf("could not execute instruction")
		}
	}
	return nil
}

func (icm *IntcodeMachine) String() string {
	return icm.Program.String()
}
