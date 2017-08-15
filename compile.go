package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const DAT = 1000
const EMPTY = 1001

// Takes a command and returns the numeric opcode
func opcode(command string) int {
	opcodes := map[string]int{
		"ADD": 100,
		"SUB": 200,
		"STA": 300,
		"LDA": 500,
		"BRA": 600,
		"BRZ": 700,
		"BRP": 800,
		"INP": 901,
		"OUT": 902,
		"HLT": 0,
		"COB": 0,
		"DAT": 1000,
	}
	op, ok := opcodes[command]
	if !ok {
		return -1
	}
	return op
}

// Regular expression used in line processing
var reg = regexp.MustCompile(`^\s*(\w*)\s*(\w*)\s*(\w*)\s*(?://)?.*$`)

// Seperates the given input into its label opcode and usedlabel
func processLine(input string) (string, int, string) {
	res := reg.FindAllStringSubmatch(input, -1)
	if len(res) == 0 {
		return "", -1, ""
	}
	if len(res[0]) == 0 {
		return "", -1, ""
	}
	line := res[0][1:]

	if line[2] == "" {
		line = line[:2]
	}
	if line[1] == "" {
		line = line[:1]
	}
	if line[0] == "" && len(line) == 1 {
		return "", EMPTY, ""
	}

	opc := -1
	lab := ""
	usedlab := ""

	opc = opcode(line[0])

	if opc == -1 {
		if len(line) == 1 {
			return "", -1, ""
		}
		opc = opcode(line[1])
		lab = line[0]
		if opc == -1 {
			return "", -1, ""
		}
		if len(line) > 2 {
			usedlab = line[2]
		}
	} else if len(line) > 1 {
		usedlab = line[1]
	}

	return lab, opc, usedlab
}

// Compiles and returns the given program
func Compile(code string) ([100]int, error) {
	var pc int
	var output [100]int

	scanner := bufio.NewScanner(strings.NewReader(code))

	labels := make(map[string]int)
	usedlabel := make(map[int]string)

	lineCount := 0

	for scanner.Scan() {
		lineCount++

		if pc > 99 {
			pc = 0
		}

		lab, opc, usedlab := processLine(scanner.Text())
		if opc == -1 {
			return output, fmt.Errorf("Error: invalid opcode on line %d", lineCount)
		}

		if opc == DAT {
			i, _ := strconv.Atoi(usedlab)
			output[pc] = wrap(i)
		} else if opc == EMPTY {
			continue
		} else {
			output[pc] = opc
			if usedlab != "" {
				if i, err := strconv.Atoi(usedlab); err == nil {
					if i < 0 || i > 99 {
						usedlabel[pc] = usedlab
					} else {
						output[pc] += i
					}
				} else {
					usedlabel[pc] = usedlab
				}
			}
		}

		if lab != "" {
			_, ok := labels[lab]
			if ok {
				return output, fmt.Errorf("Error: label %s already used", lab)
			}
			labels[lab] = pc
		}

		pc++
	}

	for key, val := range usedlabel {
		i, ok := labels[val]
		if !ok {
			return output, fmt.Errorf("Error: label %s not found", val)
		}
		output[key] += i
	}

	return output, nil
}
