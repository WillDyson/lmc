package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func executeCase(t *testing.T, dir, cas string, prog [100]int) {
	t.Log("Testing case", cas, "of example", dir)

	input, err := os.Open(path.Join(dir, "input", cas))
	if err != nil {
		t.Log("Error loading input", cas, "of example", dir)
		return
	}
	defer input.Close()

	output, err := os.Open(path.Join(dir, "output", cas))
	if err != nil {
		t.Log("Error loading output", cas, "of example", dir)
		return
	}
	defer output.Close()

	outBuf := new(bytes.Buffer)
	comp := NewComputer(input, outBuf)

	comp.Load(prog)

	err = nil
	for err == nil {
		err = comp.Step()
	}

	if err != ErrMachineOff {
		t.Error("Machine gave error for case", cas, "of example", dir)
		return
	}

	equal, err := readersEqual(output, outBuf)
	if err != nil {
		t.Error("Failed check of output for case", cas, "of example", dir)
		return
	}

	if !equal {
		t.Error("Outputs dont match for case", cas, "of example", dir)
		return
	}

	t.Log("Case", cas, "of example", dir, "successful")
}

func readersEqual(r1 io.Reader, r2 io.Reader) (bool, error) {
	b1, b2 := make([]byte, 1), make([]byte, 1)
	for {
		_, err1 := r1.Read(b1)
		_, err2 := r2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else {
				if err1 != nil {
					return false, err1
				}
				return false, err2
			}
		}

		if !bytes.Equal(b1, b2) {
			return false, nil
		}
	}
}

func executeExample(t *testing.T, dir string) {
	t.Log("Testing example", dir)

	F, err := ioutil.ReadDir(path.Join(dir, "input"))
	if err != nil {
		t.Log("Error loading inputs of example", dir)
		return
	}

	code, err := os.Open(path.Join(dir, "program.lmc"))
	if err != nil {
		t.Log("Error loading program of example", dir)
		return
	}
	defer code.Close()

	data, err := ioutil.ReadAll(code)
	if err != nil {
		t.Log("Error reading program of example", dir)
		return
	}

	prog := Compile(string(data))

	for _, file := range F {
		if file.IsDir() {
			continue
		}

		executeCase(t, dir, file.Name(), prog)
	}
}

func TestExamples(t *testing.T) {
	F, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal("Error loading examples")
	}

	for _, dir := range F {
		if !dir.IsDir() {
			continue
		}

		executeExample(t, path.Join("testdata", dir.Name()))
	}
}
