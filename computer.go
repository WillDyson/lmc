package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

// Specifies all the state of a computer
type Computer struct {
	Accumulator int
	PC          int
	Halted      bool

	Memory [100]int

	Input    io.Reader
	Output   io.Writer
	bufInput *bufio.Reader
}

// Wraps the given int into [0, 999]
func wrap(i int) int {
	i += 999
	r := i % 1999
	if r < 0 {
		r += 1999
	}
	return r - 999
}

// Separates the given int into the opcode and operand
func separate(c int) (int, int) {
	if c < 0 || c > 999 {
		return -1, 0
	}
	opcode := c / 100
	return opcode, c - opcode*100
}

var (
	ErrMachineOff     = errors.New("Machine halted")
	ErrPCOutOfBounds  = errors.New("Program counter out of bounds")
	ErrInvalidCommand = errors.New("Invalid command passed to CPU")
	ErrInputOutput    = errors.New("Error inputting")
)

// Makes a step over the loaded program
func (c *Computer) Step() error {
	if c.Halted {
		return ErrMachineOff
	}

	if c.PC < 0 || c.PC > 99 {
		return ErrPCOutOfBounds
	}

	opcode, oprand := separate(c.Memory[c.PC])

	switch opcode {
	case 1: // ADD
		c.Accumulator += c.Memory[oprand]
	case 2: // SUBTRACT
		c.Accumulator -= c.Memory[oprand]
	case 3: // STORE
		c.Memory[oprand] = c.Accumulator
	case 5: // LOAD
		c.Accumulator = c.Memory[oprand]
	case 6: // BRANCH
		c.PC = oprand - 1
	case 7: // BRANCH IF ZERO
		if c.Accumulator == 0 {
			c.PC = oprand - 1
		}
	case 8: // BRANCH IF POSITIVE
		if c.Accumulator >= 0 {
			c.PC = oprand - 1
		}
	case 9: // INPUT OUTPUT
		switch oprand {
		case 1: // INPUT
			str, err := c.bufInput.ReadString('\n')
			if len(str) == 0 || err != nil {
				return ErrInputOutput
			}
			i, err := strconv.Atoi(str[:len(str)-1])
			if err != nil {
				return ErrInputOutput
			}
			c.Accumulator = i
		case 2: // OUTPUT
			c.Output.Write([]byte(strconv.Itoa(c.Accumulator) + "\n"))
		}
	case 0: // HALT
		c.Halted = true
	default:
		return ErrInvalidCommand
	}

	c.PC++
	if c.PC > 99 {
		c.PC = 0
	}

	c.Memory[oprand] = wrap(c.Memory[oprand])
	c.Accumulator = wrap(c.Accumulator)

	return nil
}

// Loads the given program into memory and resets registers
func (c *Computer) Load(prog [100]int) error {
	for i, v := range prog {
		c.Memory[i] = wrap(v)
	}
	c.Accumulator = 0
	c.PC = 0
	return nil
}

// Creates and returns a new computer
func NewComputer(in io.Reader, out io.Writer) *Computer {
	c := &Computer{
		Input:  in,
		Output: out,
	}
	c.bufInput = bufio.NewReader(c.Input)
	return c
}
