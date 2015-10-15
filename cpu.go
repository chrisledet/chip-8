package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	StartAddress   int = 0x200
	MaxProgramSize int = 0xFFF - StartAddress
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
			os.Exit(0)
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
			fmt.Printf("CMD: subroutine return to address 0x%X\n", cpu.pc)
		} else {
			return errors.New("invalid subroutine return")
		}
	}

	switch opscode & 0xF000 {
	case 0xA000:
		cpu.i = opsval
		fmt.Printf("CMD: set register I to 0x%X\n", cpu.i)
	case 0x1000:
		cpu.pc = opsval
		return nil
	case 0x2000: // Calls subroutine at NNN.
		cpu.stack[cpu.sp] = cpu.pc
		cpu.sp++
		cpu.pc = opsval - 2 // reduce by 2 because we increment PC anyway
		fmt.Printf("CMD: jump to subroutine at address 0x%X, 0x%X\n", opsval)
		return nil
	case 0x3000:
		x := cpu.memory[pc] & 0x0F
		nn := cpu.memory[pc+1]

		fmt.Printf("CMD: skip next instruction if V%d (0x%X) = 0x%X\n", x, cpu.v[x], nn)
		if cpu.v[x] == byte(nn) {
			cpu.pc += 2
		}
	case 0x4000:
		x := cpu.memory[cpu.pc] & 0x0F
		nn := cpu.memory[cpu.pc+1]

		fmt.Printf("CMD: skip next instruction if V%d (0x%X) != 0x%X\n", x, cpu.v[x], nn)
		if cpu.v[x] != byte(nn) {
			cpu.pc += 2
		}
	case 0x5000:
		x := cpu.memory[cpu.pc] & 0x0F
		y := cpu.memory[cpu.pc+1] / 0x10

		fmt.Printf("CMD: skip next instruction if V%d (0x%X) != V%d (0x%X)\n", x, cpu.v[x], y, cpu.v[y])
		if cpu.v[x] == cpu.v[y] {
			cpu.pc += 2
		}
	case 0x6000:
		x := cpu.memory[cpu.pc] & 0x0F
		opsval = cpu.memory[cpu.pc+1]
		cpu.v[x] = byte(opsval)

		fmt.Printf("CMD: set V%d to 0x%X\n", x, opsval)
	case 0x7000:
		x := cpu.memory[cpu.pc] & 0x0F
		opsval = cpu.memory[cpu.pc+1]
		cpu.v[x] += byte(opsval)

		fmt.Printf("CMD: add 0x%X to V%d\n", opsval, x)
	case 0x8000:
		x := cpu.memory[cpu.pc] & 0x0F
		y := cpu.memory[cpu.pc+1] / 0x10

		fmt.Printf("CMD: set V%d (0x%X) to V%d (0x%X)\n", x, cpu.v[x], y, cpu.v[y])
		cpu.v[x] = cpu.v[y]
	}

	return nil
}

func NewCPU() *CPU {
	return &CPU{
		v:      make([]byte, 16),
		i:      0x0,
		pc:     0x200,
		sp:     0x0,
		memory: make([]uint16, 4096),
		stack:  make([]uint16, 16),
	}
}
