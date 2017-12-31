package machine

import (
	"bytes"
	"testing"
)

func TestLoadProgram(t *testing.T) {
	SAMPLE_BINARY := []byte{0x09, 0x00, 0x00, 0x80, 0x01, 0x80,
		0x04, 0x00, 0x13, 0x00, 0x00, 0x80}
	expected := [...]uint16{9, 32768, 32769, 4, 19, 32768}

	m := NewMachine()
	r := bytes.NewReader(SAMPLE_BINARY)
	nword, err := m.LoadProgram(r)
	if err != nil {
		t.Errorf("Error in LoadProgram(): %s\n", err)
	}
	if nword != len(expected) {
		t.Errorf("Loaded %d words (expected %d)\n", nword, len(expected))
	}

	for idx, val := range expected {
		if m.mem[idx] != val {
			t.Errorf("mem[%d] = %d (expected %d)\n", idx, m.mem[idx], val)
		}
	}
}

func TestRegisters(t *testing.T) {
	m := NewMachine()
	for regnum := uint16(0); regnum < NREG; regnum += 1 {
		value := uint16(37 - regnum)
		err := m.SetRegister(regnum, value)
		if err != nil {
			t.Errorf("Error in m.SetRegister(%d, %d): %s\n", regnum, value, err)
		}
		result, err := m.GetRegister(regnum)
		if err != nil {
			t.Errorf("Error in m.GetRegister(%d): %s\n", regnum, err)
		}
		if result != value {
			t.Errorf("m.GetRegister(%d) returned %d (expected %d)\n",
				regnum, result, value)
		}

		result, err = m.value(MAXVALUE + regnum)
		if err != nil {
			t.Errorf("Error in m.value(%d): %s\n", MAXVALUE+regnum, err)
		}
		if result != value {
			t.Errorf("m.value(%d) returned %d (expected %d)\n",
				MAXVALUE+regnum, result, value)
		}
	}
}

func TestStack(t *testing.T) {
	m := NewMachine()
	values := [...]uint16{3, 2, 1}
	nval := len(values)

	_, err := m.Pop()
	if err == nil {
		t.Errorf("m.Pop() should return error if stack is empty\n")
	}

	for _, value := range values {
		m.Push(value)
	}
	if m.sp != uint16(nval) {
		t.Errorf("m.sp is %d (expected %d)\n", m.sp, nval)
	}

	result, err := m.Pop()
	if err != nil {
		t.Errorf("Error in m.Pop(): %s\n", err)
	}
	if result != values[nval-1] {
		t.Errorf("m.Pop() returned %d (expected %d)\n", result, values[nval-1])
	}
	if m.sp != uint16(nval-1) {
		t.Errorf("m.sp is %d (expected %d)\n", m.sp, nval-1)
	}
}

func TestAdd(t *testing.T) {
	program := [...]uint16{
		uint16(ADD), MAXVALUE + REG0, 1, 1,
		uint16(ADD), MAXVALUE + REG1, MAXVALUE + REG0, 3,
		uint16(ADD), MAXVALUE + REG2, MAXVALUE + REG0, MAXVALUE + REG1,
		uint16(ADD), MAXVALUE + REG3, MAXVALUE - 1, MAXVALUE - 2,
		uint16(ADD), MAXVALUE + REG4, MAXVALUE - 7, 7}

		expected_values := [...]uint16{2, 5, 7, MAXVALUE - 3, 0}

	m := NewMachine()
	for addr, word := range program {
		m.mem[addr] = word
	}

	err := m.Execute()
	if err != nil {
		t.Errorf("Error in Execute(): %s\n", err)
	}

	for regnum, expected := range expected_values {
		result, err := m.GetRegister(uint16(regnum))
		if err != nil {
			t.Errorf("Error in m.GetRegister(%d): %s\n", regnum, err)
		}
		if result != expected {
			t.Errorf("Register %d value = %d (expected %d)\n",
				regnum, result, expected)
		}
	}
}

func TestMath (t *testing.T) {
	program := [...]uint16{
		uint16(MULT), MAXVALUE + REG0, 2, 3,
		uint16(MULT), MAXVALUE + REG1, MAXVALUE - 3, MAXVALUE - 5,
		uint16(ADD), MAXVALUE + REG2, 4, MAXVALUE - 3,
		uint16(MOD), MAXVALUE + REG3, 63, 16,
		uint16(AND), MAXVALUE + REG4, 0x3f, 0xfc,
		uint16(OR), MAXVALUE + REG5, 0x3f, 0xfc,
		uint16(NOT), MAXVALUE + REG6, MAXVALUE - 2,
		uint16(NOT), MAXVALUE + REG7, 0x5}

		expected_values := [...]uint16{6, 15, 1, 15, 0x3c, 0xff, 1, 0x7ffa}

	m := NewMachine()
	for addr, word := range program {
		m.mem[addr] = word
	}

	err := m.Execute()
	if err != nil {
		t.Errorf("Error in Execute(): %s\n", err)
	}

	for regnum, expected := range expected_values {
		result, err := m.GetRegister(uint16(regnum))
		if err != nil {
			t.Errorf("Error in m.GetRegister(%d): %s\n", regnum, err)
		}
		if result != expected {
			t.Errorf("Register %d value = %d (expected %d)\n",
				regnum, result, expected)
		}
	}
}

