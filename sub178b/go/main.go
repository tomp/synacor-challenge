package main

import (
	"flag"
	"fmt"
)

const (
	STACK_SIZE int    = 1024 * 1024 * 1024
	MODULUS    uint16 = 32768
	SOURCEFILE string = "challenge.asm"
	TRACEFILE  string = "trace.txt"
)

type Stack struct {
	stack []uint16
	sp    int
}

func confirm3(r0, r1, r7 uint) (uint16, uint16) {
	var stack [STACK_SIZE]uint16
	sp := 0
	reg0 := uint16(r0)
	reg1 := uint16(r1)
	reg7 := uint16(r7)

	var sub178b func()

	sub178b = func() {
		iter := 0
		for {
			iter += 1
			if reg0 == 0 {
				break
			} else if reg0 == 1 {
				reg1 = (reg1 + reg7) & 0x7fff
				break
			} else if reg0 == 2 {
				reg1 = (reg1 + (reg1+2)*reg7) & 0x7fff
				break
			} else if reg1 == 0 {
				reg0 -= 1
				reg1 = reg7
			} else {
				stack[sp] = reg0
				sp += 1
				reg1 -= 1
				sub178b()
				reg1 = reg0
				sp -= 1
				reg0 = stack[sp] - 1
			}
		}
		reg0 = (reg1 + 1) & 0x7fff
	}

	sub178b()
	return reg0, reg1
}

// confirm2 implements the confirmation algorithm with tail calls eliminated
func confirm2(r0, r1, r7 uint) (uint16, uint16) {
	var stack [STACK_SIZE]uint16
	sp := 0
	reg0 := uint16(r0)
	reg1 := uint16(r1)
	reg7 := uint16(r7)

	var sub178b func()

	sub178b = func() {
		iter := 0
		for {
			iter += 1
			if reg0 == 0 {
				reg0 = (reg1 + 1) & 0x7fff
				break
			} else if reg1 == 0 {
				reg0 -= 1
				reg1 = reg7
			} else {
				stack[sp] = reg0
				sp += 1
				reg1 -= 1
				sub178b()
				reg1 = reg0
				sp -= 1
				reg0 = stack[sp] - 1
			}
		}
	}

	sub178b()
	return reg0, reg1
}

// confirm() is a direct translation of the machine-language confirmation
// algorithm to Go, with a fixed-size stack.
func confirm(r0, r1, r7 uint) (uint16, uint16) {
	var stack [STACK_SIZE]uint16
	sp := 0
	reg0 := uint16(r0)
	reg1 := uint16(r1)
	reg7 := uint16(r7)

	var sub178b func()

	sub178b = func() {
		if reg0 == 0 {
			reg0 = (reg1 + 1) & 0x7fff
		} else if reg1 == 0 {
			reg0 -= 1
			reg1 = reg7
			sub178b()
		} else {
			stack[sp] = reg0
			sp += 1
			reg1 -= 1
			sub178b()
			reg1 = reg0
			sp -= 1
			reg0 = stack[sp] - 1
			sub178b()
		}
		return
	}

	sub178b()
	return reg0, reg1
}

func main() {
	fmt.Println("teleporter confirmation code")
	var r0 = flag.Uint("r0", 0, "initial value for register 0")
	var r1 = flag.Uint("r1", 0, "initial value for register 1")
	var r7 = flag.Uint("r7", 0, "initial value for register 7")
	var v2 = flag.Bool("v2", false, "use version 2 of the confirmation code")
	var v3 = flag.Bool("v3", false, "use version 3 of the confirmation code")
	var solve = flag.Bool("solve", false, "find the teleporter code")
	flag.Parse()

	var reg0, reg1 uint16
	var reg7 uint

	if *solve {
		for reg7 = 1; reg7 < 32767; reg7 += 1 {
			reg0, reg1 = confirm3(4, 1, reg7)
			if reg0 == 6 {
				fmt.Printf("(%d, %d, %d) -> r0: %d  r1: %d\n",
					4, 1, reg7, reg0, reg1)
			}
		}
	} else {
		if *v3 {
			reg0, reg1 = confirm3(*r0, *r1, *r7)
		} else if *v2 {
			reg0, reg1 = confirm2(*r0, *r1, *r7)
		} else {
			reg0, reg1 = confirm(*r0, *r1, *r7)
		}
		fmt.Printf("(%d, %d, %d) -> r0: %d  r1: %d\n", *r0, *r1, *r7, reg0, reg1)
	}

}
