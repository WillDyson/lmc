package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"os"
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

var errNoFileGiven = errors.New("No argument for file was given")

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

func compile(c *cli.Context) error {
	var reader io.Reader

	if appTerminal {
		reader = os.Stdin
	} else {
		if c.NArg() == 0 {
			return errNoFileGiven
		}

		file, err := os.Open(c.Args().Get(0))
		if err != nil {
			fmt.Println(err)
			return err
		}
		reader = file
		defer file.Close()
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = writeOutput(Compile(string(data)), appOut)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Prints useful information about the current step to the user
func printStepInfo(c *Computer) {
	fmt.Println("PC:", c.PC, "ACC:", c.Accumulator, "CODE:", c.Memory[c.PC])
}

func run(c *cli.Context) error {
	if c.NArg() == 0 {
		return errNoFileGiven
	}

	prog, err := readOutput(c.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		return err
	}

	comp := NewComputer(os.Stdin, os.Stdout)
	comp.Load(prog)

	err = nil
	for err == nil {
		printStepInfo(comp)
		err = comp.Step()
	}

	fmt.Println(err)

	if err == ErrMachineOff {
		return nil
	}

	return err
}

var (
	appTerminal bool
	appFile     string
	appOut      string
)

func main() {
	app := cli.NewApp()
	app.Name = "LMC"
	app.Usage = "compile and execute little man programs"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "compile",
			Aliases: []string{"c"},
			Usage:   "compile and write the given program",
			Action:  compile,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "terminal,t",
					Usage:       "input program using stdin until EOF",
					Destination: &appTerminal,
				},
				cli.StringFlag{
					Name:        "out,o",
					Usage:       "specifies the file to be written",
					Destination: &appOut,
					Value:       "a.out",
				},
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "runs the given program",
			Action:  run,
		},
	}
	app.Run(os.Args)
}
