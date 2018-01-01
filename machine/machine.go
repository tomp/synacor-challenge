package machine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type valType uint16
type opCode uint16
type regType uint16
type argType uint16

type opType struct {
	code opCode  // the instruction code
	name string  // the instruction name (for disassembly)
	args argType // number and type of instruction arguments
}

const NREG = 8          // number of registers
const MAXVALUE = 0x8000 // upper bound on 15-bit values
const BITMASK = 0x7fff  // mask for 15-bit values

const MEMSIZE = MAXVALUE // size of main memory
const STACKSIZE = 1024   // initial stack size

// Code indictating the number and type of an instruction's arguments
const (
	Z   argType = iota
	V           // value
	VV          // value, value
	R           // register
	RV          // register, value
	RVV         // register, value, value
	C // character
	A // address
	AV // address, value
	VA // value, address
	RA // register, address
)

// Constants for instruction opcodes.
const (
	HALT opCode = iota
	SET
	PUSH
	POP
	EQ
	GT
	JMP
	JT
	JF
	ADD
	MULT
	MOD
	AND
	OR
	NOT
	RMEM
	WMEM
	CALL
	RET
	OUT
	IN
	NOOP
)

// Constants for register numbers, to make test code more readable.
const (
	REG0 uint16 = iota
	REG1
	REG2
	REG3
	REG4
	REG5
	REG6
	REG7
)

// This is a slice, not a map, and values are looked up using the
// opcode as the index.  So, each instruction's data needs to be at
// the correct locaction, so that OPS[idx].code == idx.
var OPS = [...]opType{
	{HALT, "halt", Z},
	{SET, "set", RV},
	{PUSH, "push", V},
	{POP, "pop", R},
	{EQ, "eq", RVV},
	{GT, "gt", RVV},
	{JMP, "jmp", A},
	{JT, "jt", VA},
	{JF, "jf", VA},
	{ADD, "add", RVV},
	{MULT, "mult", RVV},
	{MOD, "mod", RVV},
	{AND, "and", RVV},
	{OR, "or", RVV},
	{NOT, "not", RV},
	{RMEM, "rmem", RA},
	{WMEM, "wmem", AV},
	{CALL, "call", A},
	{RET, "ret", Z},
	{OUT, "out", C},
	{IN, "in", R},
	{NOOP, "noop", Z},
}

type Machine struct {
	mem   [MEMSIZE]uint16 // memory
	reg   [NREG]uint16    // registers
	stack []uint16        // stack
	pc    uint16          // program counter (mem addr)
	sp    uint16          // stack pointer
	rdr   *bufio.Reader   // buffered reader for stdin
}

func NewMachine() (m *Machine) {
	m = new(Machine)
	m.stack = make([]uint16, STACKSIZE)
	m.pc = 0
	m.sp = 0
	m.rdr = bufio.NewReader(os.Stdin)
	return
}

// ReadUint16 reads two bytes and interprets them as a little-endian Uint16
// value.  If two bytes aren't available, returns an error.
func readUint16(b *bufio.Reader) (uint16, error) {
	b1, err := b.ReadByte()
	if err == nil {
		b2, err := b.ReadByte()
		if err == nil {
			return uint16(b1) + uint16(b2)<<8, nil
		}
	}
	return 0, err
}

// LoadProgram reads bytes from the given io.Reader, interprets them as
// little-endian Uint16 values, and stores them sequentially in memory,
// starting at address zero.
// The number of 16-bit words loaded is returned, along with an error code.
func (m *Machine) LoadProgram(r io.Reader) (int, error) {
	var err error
	var value uint16
	addr := 0
	rdr := bufio.NewReader(r)
	for err == nil {
		value, err = readUint16(rdr)
		if err == nil {
			m.mem[addr] = value
			addr += 1
		}
	}
	if err != io.EOF {
		return addr, err
	}
	return addr, nil
}

// RegisterNumber returns the register number for the given number.
func RegisterNumber(number uint16) (uint16, error) {
	if number < MAXVALUE || number > MAXVALUE+NREG {
		return 0, fmt.Errorf("Invalid register reference: %0x", number)
	}
	return number & 0x0007, nil
}

// GetRegister returns the value of the specified register
func (m *Machine) GetRegister(regnum uint16) (uint16, error) {
	if regnum < NREG {
		return m.reg[regnum], nil
	}
	return 0, fmt.Errorf("No such register: %d", regnum)
}

// SetRegister assigns a new value to the specified register
func (m *Machine) SetRegister(regnum uint16, value uint16) error {
	if regnum < NREG {
		m.reg[regnum] = value
		return nil
	}
	return fmt.Errorf("No such register: %d", regnum)
}

// Push pushes the given value to the stack.
func (m *Machine) Push(value uint16) {
	if m.sp == uint16(len(m.stack)) {
		m.stack = append(m.stack, value)
	} else {
		m.stack[m.sp] = value
	}
	m.sp += 1
}

// Pop removes the top value from the stack and returns it.
// An error is returned if the stack is empty.
func (m *Machine) Pop() (uint16, error) {
	if m.sp == 0 {
		return 0, fmt.Errorf("Stack is empty")
	}
	m.sp -= 1
	return m.stack[m.sp], nil
}

// value returns the value represented by the given number.
func (m *Machine) value(number uint16) (uint16, error) {
	if number < MAXVALUE {
		return number, nil
	}
	regnum, err := RegisterNumber(number)
	if err != nil {
		return 0, err
	}
	return m.reg[regnum], nil
}

// GetInstruction returns the opcode at the given memory location,
// and its arguments.  (There may be 1, 2, 3 or no arguments,
// depending on the instruction.)
func (m *Machine) GetInstruction(pc uint16) (op opCode, a, b, c, next_pc uint16, err error) {
	if opCode(m.mem[pc]) > NOOP {
		err = fmt.Errorf("Illegal opcode %d at pc=0x%04x\n", m.mem[pc], pc)
		return
	}
	op = opCode(m.mem[pc])
	switch OPS[op].args {
	case R:
		a, err = RegisterNumber(m.mem[pc+1])
		b, c = 0, 0
		next_pc = pc + 2
	case RV, RA:
		a, err = RegisterNumber(m.mem[pc+1])
		if err == nil {
			b, err = m.value(m.mem[pc+2])
		}
		c = 0
		next_pc = pc + 3
	case RVV:
		a, err = RegisterNumber(m.mem[pc+1])
		if err == nil {
			b, err = m.value(m.mem[pc+2])
		}
		if err == nil {
			c, err = m.value(m.mem[pc+3])
		}
		next_pc = pc + 4
	case V, A, C:
		a, err = m.value(m.mem[pc+1])
		b, c = 0, 0
		next_pc = pc + 2
	case VV, VA, AV:
		a, err = m.value(m.mem[pc+1])
		if err == nil {
			b, err = m.value(m.mem[pc+2])
		}
		c = 0
		next_pc = pc + 3
	case Z:
		next_pc = pc + 1
	default:
		err = fmt.Errorf("Unrecognized args type: %d", OPS[op].args)
	}
	return
}

func FormatVal(value uint16) string {
	regnum, err := RegisterNumber(value)
	if err == nil {
		return fmt.Sprintf("reg%d", regnum)
	} else if value > 0xff {
		return fmt.Sprintf("%04x", value)
	} else {
		return fmt.Sprintf("%d", value)
	}
}

func FormatWords(values []uint16, nword int) string {
	fields := make([]string, nword)
	for idx, value := range(values) {
		fields[idx] = fmt.Sprintf("0x%04x", value)
	}
	return strings.Join(fields, "  ")

}

// FormatInstruction returns a string representing the instruction in
// "disassembled" form.  It is assumed that the arguments come from a prior
// call to GetInstruction, and therefore don't need to be validated again.
func FormatInstruction(words []uint16) string {
	op := opCode(words[0])
	name := OPS[op].name
	switch OPS[op].args {
	case R, V:
		return fmt.Sprintf("%s %s", name, FormatVal(words[1]))
	case C:
		return fmt.Sprintf("%s %q", name, words[1])
	case A:
		return fmt.Sprintf("%s @%s", name, FormatVal(words[1]))
	case RV, VV:
		return fmt.Sprintf("%s %s, %s", name, FormatVal(words[1]),
		FormatVal(words[2]))
	case RA:
		return fmt.Sprintf("%s %s, @%s", name, FormatVal(words[1]),
		FormatVal(words[2]))
	case RVV:
		return fmt.Sprintf("%s %s, %s, %s", name, FormatVal(words[1]),
		FormatVal(words[2]), FormatVal(words[3]))
	case AV:
		return fmt.Sprintf("%s @%s, %s", name, FormatVal(words[1]),
		FormatVal(words[2]))
	case VA:
		return fmt.Sprintf("%s %s, @%s", name, FormatVal(words[1]),
		FormatVal(words[2]))
	default:
		return name
	}
}

// Execute the program starting at the current program counter.
// An error is returned if execution doesn't terminate on a HALT
// instruction.
func (m *Machine) Execute() (err error) {
	var val uint16
	var ch byte
	for err == nil {
		op, a, b, c, next_pc, err := m.GetInstruction(m.pc)
		if err != nil {
			break
		}

		switch op {
		case HALT:
			return nil
		case SET:
			err = m.SetRegister(a, b)
		case PUSH:
			m.Push(a)
		case POP:
			val, err = m.Pop()
			if err == nil {
				err = m.SetRegister(a, val)
			}
		case EQ:
			if b == c {
				err = m.SetRegister(a, 1)
			} else {
				err = m.SetRegister(a, 0)
			}
		case GT:
			if b > c {
				err = m.SetRegister(a, 1)
			} else {
				err = m.SetRegister(a, 0)
			}
		case JMP:
			next_pc = a
		case JT:
			if a != 0 {
				next_pc = b
			}
		case JF:
			if a == 0 {
				next_pc = b
			}
		case ADD:
			err = m.SetRegister(a, (b+c)&BITMASK)
		case MULT:
			// OK, because operands are uint16
			err = m.SetRegister(a, (b*c)&BITMASK)
		case MOD:
			err = m.SetRegister(a, b%c)
		case AND:
			err = m.SetRegister(a, b&c)
		case OR:
			err = m.SetRegister(a, b|c)
		case NOT:
			err = m.SetRegister(a, BITMASK^b)
		case RMEM:
			err = m.SetRegister(a, m.mem[b])
		case WMEM:
			m.mem[a] = b
		case CALL:
			m.Push(next_pc)
			next_pc = a
		case RET:
			next_pc, err = m.Pop()
		case OUT:
			fmt.Printf("%c", a)
		case IN:
			ch, err = m.rdr.ReadByte()
			if err == nil {
				err = m.SetRegister(a, uint16(ch))
			}
		case NOOP:
		default:
			err = fmt.Errorf("Unrecognized opcode: %d at pc=%d", op, m.pc)
		}
		if err == nil {
			m.pc = next_pc
		}
	}
	return err
}

func (m *Machine) Disassemble(pc int) error {
	addr_queue := []uint16{uint16(pc)}
	visited := make(map[uint16]bool)
	for len(addr_queue) > 0 {
		addr := addr_queue[0]
		addr_queue = addr_queue[1:]
		for addr < MAXVALUE {
			if visited[addr] {
				break
			}
			visited[addr] = true
			op, _, _, _, next_pc, err := m.GetInstruction(addr)
			if err != nil {
				break
			}
			values := m.mem[addr:next_pc]
			words := FormatWords(values, 4)
			code := FormatInstruction(values)
			fmt.Printf("%04x: %-30s : %s\n", addr, words, code)
			for _, value := range(values[1:]) {
				if value < MAXVALUE && value > 0xff && ! visited[value] {
					addr_queue = append(addr_queue, value)
				}
			}
			if op == JMP || op == HALT {
				addr = MAXVALUE
			} else {
				addr = next_pc
			}
		}
	}
	return nil
}
