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
	ErrBadMemoryAccess      MachineError = `unexpected access to memory address`
	EOX                     MachineError = `end of execution`
)

type MachineError string

func (e MachineError) Error() string {
	return string(e)
}

type IntcodeMachine struct {
	Program Program
	ISet    InstructionSet
	rawData []byte
	sep     string
}

func parseLine(symbols []string) (Instruction, int, error) {
	op, err := strconv.Atoi(symbols[0])
	if err != nil {
		return Instruction{}, 0, fmt.Errorf("could not parse instruction '%s' : %w", symbols[0], ErrMalformedInstruction)
	}
	if len(symbols) == 1 {
		if symbols[0] == HaltSym {
			return Instruction{
				op,
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
		op,
		term1,
		term2,
		store,
	}, 4, nil
}

func parseProgram(data []byte, sep string) (Program, error) {
	program := Program{}
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
	return program, nil
}
func NewIntcodeMachine(r io.Reader, sep string) (*IntcodeMachine, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate machine: %w", err)
	}
	program, err := parseProgram(data, sep)
	if err != nil {
		return nil, fmt.Errorf("could not parse program: %w", err)
	}
	return &IntcodeMachine{
		Program: program,
		ISet:    make(InstructionSet),
		rawData: data,
		sep:     sep,
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

func (icm *IntcodeMachine) AddInstruction(code OpCode, fn func(int, int) (int, error)) error {
	icm.ISet[code] = fn
	return nil
}

func (icm *IntcodeMachine) executeInstruction(i Instruction, ip int) error {
	var (
		a, b int
		err  error
	)
	fn, ok := icm.ISet[i[0]]
	if !ok {
		return fmt.Errorf("could not execute instruction %s: %w", i, ErrUnknownOpCode)
	}
	if i[0] == HaltOp {
		a, b = 0, 0
	} else {
		a = icm.Program[i[1]]
		b = icm.Program[i[2]]
	}
	v, err := fn(a, b)
	if err != nil {
		if err == EOX {
			return EOX
		}
		return fmt.Errorf("could not execute instruction %s", i)
	}
	icm.Program[i[3]] = v
	return nil
}
func (icm *IntcodeMachine) fetch(ip int) (Instruction, int) {
	var (
		instr Instruction
		size  int
	)
	if icm.Program[ip] == HaltOp {
		instr = Instruction{icm.Program[ip]}
		size = 1
	} else {
		instr = Instruction{
			icm.Program[ip],
			icm.Program[ip+1],
			icm.Program[ip+2],
			icm.Program[ip+3],
		}
		size = 4
	}
	return instr, size
}
func (icm *IntcodeMachine) Execute() error {
	for ip := 0; ip < len(icm.Program); {
		instr, size := icm.fetch(ip)
		err := icm.executeInstruction(instr, ip)
		if err != nil {
			if err == EOX {
				return EOX
			}
			return fmt.Errorf("could not execute instruction")
		}
		ip += size
	}
	return nil
}

func (icm *IntcodeMachine) Print(addr int) (string, error) {
	if addr >= len(icm.Program) {
		return "", fmt.Errorf("could not set code: %w", ErrBadMemoryAccess)
	}
	return fmt.Sprintf("%d", icm.Program[addr]), nil
}

func (icm *IntcodeMachine) Reset() {
	program, _ := parseProgram(icm.rawData, icm.sep)
	icm.Program = program
}

func (icm *IntcodeMachine) SetOpCode(addr int, code OpCode) error {
	if addr >= len(icm.Program) {
		return fmt.Errorf("could not set code: %w", ErrBadMemoryAccess)
	}
	icm.Program[addr] = code
	return nil
}

func (icm *IntcodeMachine) String() string {
	return icm.Program.String()
}
