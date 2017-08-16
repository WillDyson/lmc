package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"os"
)

var errNoFileGiven = errors.New("No argument for file was given")

func compile(c *cli.Context) error {
	var reader io.Reader

	if appTerminal {
		reader = os.Stdin
	} else {
		if c.NArg() == 0 {
			fmt.Println("No file was given to be compiled")
			return errNoFileGiven
		}

		file, err := os.Open(c.Args().Get(0))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to open program for reading")
			return err
		}
		reader = file
		defer file.Close()
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to read program")
		return err
	}

	prog, err := Compile(string(data))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to compile program")
		return err
	}

	err = writeOutput(prog, appOut)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to write the compiled program to hd")
		return err
	}

	return nil
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
		if !appQuiet {
			printStepInfo(comp)
		}
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
	appQuiet    bool
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "quiet,q",
					Usage:       "hides the machine state information",
					Destination: &appQuiet,
				},
			},
		},
	}
	app.Run(os.Args)
}
