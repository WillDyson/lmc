package main

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
)

func writeOutput(prog [100]int, file string) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(prog)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, buf.Bytes(), 0644)
}

func readOutput(file string) ([100]int, error) {
	var prog [100]int

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return prog, err
	}

	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err = dec.Decode(&prog)

	return prog, err
}
