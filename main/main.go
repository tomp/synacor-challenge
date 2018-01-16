package main

import (
	"fmt"
	"github.com/tomp/synacor-challenge/machine"
	"io"
	"os"
)

const (
	INPUTFILE  string = "challenge.bin"
	SOURCEFILE string = "challenge.asm"
	TRACEFILE  string = "trace.txt"
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

	fp, err := os.Create(SOURCEFILE)
	if err == nil {
		fmt.Printf("Writing disassembled code to %s\n", SOURCEFILE)
	}
	addrChan := make(chan uint16, 0)
	go func(w io.WriteCloser, addresses chan uint16) {
		defer w.Close()
		fmt.Fprintln(w, "================================")
		err = m.Disassemble(w, addrChan)
		fmt.Fprintln(w, "================================")
		if err != nil {
			fmt.Fprintf(w, "Error: %s\n", err)
		}
	}(fp, addrChan)

	addrChan <- uint16(0)

	fp2, err := os.Create(TRACEFILE)
	if err == nil {
		fmt.Printf("Writing trace to %s\n", TRACEFILE)
	}

	err = m.Execute(addrChan, fp2)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
