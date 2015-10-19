package main

import (
	"errors"
	"fmt"
)

const (
	StartAddress   int  = 0x200
	MaxProgramSize int  = 0xFFF - StartAddress
	CarryFlag      byte = 0xF
	RegisterCount  byte = 0x10
	StackSize      byte = 0x10
)

type CPU struct {
	v  []byte
	i  uint16
	pc uint16

	// 0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	// 0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	// 0x200-0xFFF - Program ROM and work RAM
	memory []uint16
	stack  []uint16
	sp     byte

	delayTimer uint16
	soundTimer uint16
}

func (cpu *CPU) Clear() {
	cpu.sp = 0
	for x, _ := range cpu.v {
		cpu.v[x] = 0
	}
}

func (cpu *CPU) Load(program []byte) error {
	if len(program) > MaxProgramSize {
		return errors.New("Program is too larged. Must be less than " + string(MaxProgramSize) + " bytes")
	}

	for i, opscode := range program {
		cpu.memory[StartAddress+i] = uint16(opscode)
	}

	return cpu.start()
}

func (cpu *CPU) start() error {
	for {
		err := cpu.execute()
		if err != nil {
			fmt.Printf("ERR: %s\n", err.Error())
		}

		cpu.pc += 2

		if cpu.pc > 0xFFE {
			return nil
		}
	}
}

func (cpu *CPU) execute() error {
	opscode := cpu.memory[cpu.pc]<<8 | cpu.memory[cpu.pc+1]
	opsval := opscode & 0x0FFF
	pc := cpu.pc

	if opscode == 0x00EE {
		if cpu.sp > 0 {
			cpu.pc = cpu.stack[cpu.sp-1]
			cpu.sp--
		} else {
			return errors.New("invalid subroutine return")
		}
	}

	switch opscode & 0xF000 {
	case 0xA000:
		cpu.i = opsval
	case 0x1000:
		cpu.pc = opsval
	case 0x2000: // Calls subroutine at NNN.
		// TODO: raise error if cpu.sp is 0xF
		cpu.stack[cpu.sp] = cpu.pc
		cpu.sp++
		cpu.pc = opsval - 2 // reduce by 2 because we increment PC anyway
		return nil
	case 0x3000:
		x := cpu.memory[pc] & 0x0F
		nn := cpu.memory[pc+1]

		if cpu.v[x] == byte(nn) {
			cpu.pc += 2
		}
	case 0x4000:
		x := cpu.memory[cpu.pc] & 0x0F
		nn := cpu.memory[cpu.pc+1]

		if cpu.v[x] != byte(nn) {
			cpu.pc += 2
		}
	case 0x5000:
		x := cpu.memory[cpu.pc] & 0x0F
		y := cpu.memory[cpu.pc+1] / 0x10

		if cpu.v[x] == cpu.v[y] {
			cpu.pc += 2
		}
	case 0x6000:
		x := cpu.memory[cpu.pc] & 0x0F
		opsval = cpu.memory[cpu.pc+1]
		cpu.v[x] = byte(opsval)
	case 0x7000:
		x := cpu.memory[cpu.pc] & 0x0F
		opsval = cpu.memory[cpu.pc+1]
		cpu.v[x] += byte(opsval)
	case 0x8000:
		x := cpu.memory[cpu.pc] & 0x0F
		y := cpu.memory[cpu.pc+1] / 0x10
		cond := cpu.memory[cpu.pc+1] & 0x0F

		if cond == 0x0 {
			cpu.v[x] = cpu.v[y]
		} else if cond == 0x1 {
			cpu.v[x] = cpu.v[x] | cpu.v[y]
		} else if cond == 0x2 {
			cpu.v[x] = cpu.v[x] & cpu.v[y]
		} else if cond == 0x3 {
			cpu.v[x] = cpu.v[x] ^ cpu.v[y]
		} else if cond == 0x4 {
			result := uint16(cpu.v[x]) + uint16(cpu.v[y])

			if result > 0xFF {
				cpu.v[CarryFlag] = 0x1
			} else {
				cpu.v[CarryFlag] = 0x0
			}
		} else if cond == 0x5 {
			result := int16(cpu.v[x]) - int16(cpu.v[y])

			if result < 0 {
				cpu.v[x] = byte(result + 0xFF)
				cpu.v[CarryFlag] = 0x1
			} else {
				cpu.v[x] = byte(result)
				cpu.v[CarryFlag] = 0x0
			}
		}
	}

	return nil
}

func NewCPU() *CPU {
	return &CPU{
		v:      make([]byte, RegisterCount),
		i:      0x0,
		pc:     0x200,
		sp:     0x0,
		memory: make([]uint16, 4096),
		stack:  make([]uint16, StackSize),
	}
}
