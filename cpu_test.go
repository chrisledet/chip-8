package main

import (
	"testing"
)

// Set I to NNN
func Test0xA000(t *testing.T) {
	cpu := NewCPU()
	program := []byte{0xA1, 0x23}

	cpu.Load(program)

	if cpu.i != 0x123 {
		t.Errorf("Expected 0xA000 to set I to 0x%X\n", 0x123)
	}
}

// Set VX to value of VY
func Test0x8XY0(t *testing.T) {
	cpu := NewCPU()
	program := []byte{
		0x60, 0x22, // set V0 to 0x22
		0x81, 0x00, // set V1 to V0
	}

	cpu.Load(program)

	if cpu.v[0] != 0x22 {
		t.Errorf("Expected V0 to be 0x%X but was 0x%X\n", 0x22, cpu.v[0])
	}

	if cpu.v[1] != 0x22 {
		t.Errorf("Expected V1 to be 0x%X but was 0x%X\n", 0x22, cpu.v[1])
	}
}

// Set VX to (VX OR VY)
func Test0x8XY1(t *testing.T) {
	cpu := NewCPU()
	program := []byte{
		0x60, 0x35, // set V0 to 0x35
		0x61, 0xF0, // set V1 to 0xF0
		0x80, 0x11, // V0 = 0x35 | 0xF0 (0x30)
	}

	cpu.Load(program)

	expected := byte(0xF5)
	if cpu.v[0] != expected {
		t.Errorf("Expected V0 to be 0x%X but was 0x%X\n", expected, cpu.v[0])
	}
}

// Set VX to (VX AND VY)
func Test0x8XY2(t *testing.T) {
	cpu := NewCPU()
	program := []byte{
		0x60, 0x35, // set V0 to 0x35
		0x61, 0xF0, // set V1 to 0xF0
		0x80, 0x12, // V0 = 0x35 & 0xF0 (0xF5)
	}

	cpu.Load(program)

	expected := byte(0x30)
	if cpu.v[0] != expected {
		t.Errorf("Expected V0 to be 0x%X but was 0x%X\n", expected, cpu.v[0])
	}
}

// Set VX to (VX XOR VY)
func Test0x8XY3(t *testing.T) {
	cpu := NewCPU()
	program := []byte{
		0x60, 0x35, // set V0 to 0x35
		0x61, 0x10, // set V1 to 0xF0
		0x80, 0x13, // V0 = 0x35 | 0xF0 (0x30)
	}

	cpu.Load(program)

	expected := byte(0x25)
	if cpu.v[0] != expected {
		t.Errorf("Expected V0 to be 0x%X but was 0x%X\n", expected, cpu.v[0])
	}
}

// 8XY4 - Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
func Test0x8XY4WithoutCarry(t *testing.T) {
	var expected byte
	cpu := NewCPU()
	program := []byte{
		0x60, 0x00, // set V0 to 0x00
		0x61, 0x0F, // set V1 to 0x0F
		0x80, 0x14,
	}

	cpu.Load(program)

	expected = 0x0
	if cpu.v[0] != expected {
		t.Errorf("Expected V0 to be 0x%X but was 0x%X\n", expected, cpu.v[0])
	}

	expected = 0x0F
	if cpu.v[1] != expected {
		t.Errorf("Expected V1 to be 0x%X but was 0x%X\n", expected, cpu.v[1])
	}

	expected = 0x0
	if cpu.v[15] != expected {
		t.Errorf("Expected VF to be 0x%X but was 0x%X\n", expected, cpu.v[15])
	}
}

// 8XY4 - Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
func Test0x8XY4WithCarry(t *testing.T) {
	var expected byte
	cpu := NewCPU()
	program := []byte{
		0x60, 0xFF, // set V0 to 0xFF
		0x61, 0x0F, // set V1 to 0x0F
		0x80, 0x14,
	}

	cpu.Load(program)

	expected = 0xFF
	if cpu.v[0] != expected {
		t.Errorf("Expected V0 to be 0x%X but was 0x%X\n", expected, cpu.v[0])
	}

	expected = 0x0F
	if cpu.v[1] != expected {
		t.Errorf("Expected V1 to be 0x%X but was 0x%X\n", expected, cpu.v[1])
	}

	expected = 0x1
	if cpu.v[15] != expected {
		t.Errorf("Expected VF to be 0x%X but was 0x%X\n", expected, cpu.v[15])
	}
}

// 8XY6 - Shifts VX right by one.
// 				VF is set to the value of the least significant bit of VX before the shift.
func Test0x8XY6(t *testing.T) {

}
