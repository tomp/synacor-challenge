package main

import (
	"fmt"
	"github.com/tomp/synacor-challenge/machine"
	"os"
)

const (
	INPUTFILE string = "challenge.bin"
)

func main() {
	m := machine.NewMachine()

	fmt.Printf("Load program from %s\n", INPUTFILE)
	r, err := os.Open(INPUTFILE)
	if err == nil {
		defer r.Close()
		nword, err := m.LoadProgram(r)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Loaded %d words into memory\n", nword)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("================================")
	err = m.Disassemble(0)
	fmt.Println("================================")

	// err = m.Execute()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
