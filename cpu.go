package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/chrisledet/c8vm/gfx"
)

const (
	StartAddress   int  = 0x200
	MaxProgramSize int  = 0xFFF - StartAddress
	RegisterCount  byte = 0x10
	StackSize      byte = 0x10
)

type CPU struct {
	v  []byte
	i  uint16
	pc uint16

	display Display

	// 0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	// 0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	// 0x200-0xFFF - Program ROM and work RAM
	memory []uint16
	stack  []uint16
	sp     byte

	delayTimer uint16
	soundTimer uint16

	window gfx.Window
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

	cpu.start()
	return nil
}

func (cpu *CPU) start() {
	for {
		err := cpu.execute()
		if err != nil {
			fmt.Printf("ERR: %s\n", err.Error())
		}

		cpu.pc += 2

		if cpu.pc > 0xFFE {
			break
		}
	}

	cpu.window.Close()
}

func (cpu *CPU) drawSprite(x, y, length int) {
	for h := 0; h < length; h++ {
		for w := 0; w < 8; w++ {
			fmt.Printf("cpu.memory[%d]: 0x%X\t\t", cpu.i+uint16(h), cpu.memory[cpu.i+uint16(h)])
			fmt.Printf("(0x80 >> uint16(0x%X)\n", uint16(w))

			// fmt.Printf("cpu.memory[%d]\n", cpu.i+uint16(w))
			if cpu.memory[cpu.i+uint16(h)] > (0x80 >> uint16(w)) {
				fmt.Printf("setting display[%d][%d] to true\n", x+w, y+h)
				cpu.display[x+w][y+h] = true
			}
		}
	}
}

func (cpu *CPU) execute() error {
	opscode := cpu.memory[cpu.pc]<<8 | cpu.memory[cpu.pc+1]
	opsval := opscode & 0x0FFF
	pc := cpu.pc

	switch opscode & 0xF000 {
	case 0x0000:
		switch cpu.memory[cpu.pc+1] {
		case 0xEE:
			if cpu.sp > 0 {
				cpu.pc = cpu.stack[cpu.sp-1]
				cpu.sp--
			} else {
				return errors.New("invalid subroutine return")
			}
		case 0xE0:
			cpu.display.Clear()
		}
	case 0xA000:
		cpu.i = opsval
		fmt.Printf("Setting I = 0x%X\n", opsval)
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
		fmt.Printf("Setting V%X = 0x%X\n", x, opsval)
	case 0x7000:
		x := cpu.memory[cpu.pc] & 0x0F
		opsval = cpu.memory[cpu.pc+1]
		cpu.v[x] += byte(opsval)
	case 0x8000:
		x := cpu.memory[cpu.pc] & 0x0F
		y := cpu.memory[cpu.pc+1] / 0x10
		cond := cpu.memory[cpu.pc+1] & 0x0F

		switch cond {
		case 0x0:
			cpu.v[x] = cpu.v[y]
		case 0x1:
			cpu.v[x] = cpu.v[x] | cpu.v[y]
		case 0x2:
			cpu.v[x] = cpu.v[x] & cpu.v[y]
		case 0x3:
			var result byte
			if cpu.v[x] != cpu.v[y] {
				result = cpu.v[x] - cpu.v[y]
			} else {
				result = 0x0
			}
			cpu.v[x] = result
		case 0x4:
			result := uint16(cpu.v[x]) + uint16(cpu.v[y])

			if result > 0xFF {
				cpu.v[0xF] = 0x1
			} else {
				cpu.v[0xF] = 0x0
			}
		case 0x5:
			result := int16(cpu.v[x]) - int16(cpu.v[y])

			if result < 0 {
				cpu.v[x] = byte(result + 0xFF)
				cpu.v[0xF] = 0x1
			} else {
				cpu.v[x] = byte(result)
				cpu.v[0xF] = 0x0
			}
		case 0x6:
			cpu.v[0xF] = cpu.v[x] & 0xF
			cpu.v[x] = cpu.v[x] >> 1
		case 0x7:
			cpu.v[x] = cpu.v[y] - cpu.v[x]

			if cpu.v[x] > cpu.v[y] {
				cpu.v[0xF] = 0x1
			}
		case 0xE:
			cpu.v[0xF] = cpu.v[x] / 0x10
			cpu.v[x] = cpu.v[x] << 1
		}
	case 0x9000:
		x := cpu.memory[cpu.pc] & 0x0F
		y := cpu.memory[cpu.pc+1] / 0x10

		if cpu.v[x] != cpu.v[y] {
			cpu.pc += 2
		}
	case 0xB000:
		cpu.pc = opsval + uint16(cpu.v[0x0]) - 2
	case 0xD000:
		x := cpu.memory[cpu.pc] & 0x0F
		y := cpu.memory[cpu.pc+1] / 0x10
		n := cpu.memory[cpu.pc+1] & 0x0F

		fmt.Printf("Drawing sprite at (%d, %d) for %d height\n", cpu.v[x], cpu.v[y], n)
		cpu.drawSprite(int(cpu.v[x]), int(cpu.v[y]), int(n))
		cpu.display.ToConsole()

	case 0xF000:
		x := cpu.memory[cpu.pc] & 0x0F
		nn := cpu.memory[cpu.pc+1]

		switch nn {
		case 0x07:
			// ignore
		case 0x0A:
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Waiting for Input: ")
			input, err := reader.ReadByte()

			if err != nil {
				return errors.New("error receiving input: " + err.Error())
			}

			fmt.Printf("Received: 0x%X\n", input)
			cpu.v[x] = input
		case 0x15:
			// ignore
		case 0x18:
			// ignore
		case 0x1E:
			cpu.i += uint16(cpu.v[x])
		}
	}

	return nil
}

func NewCPU() *CPU {
	cpu := CPU{
		v:      make([]byte, RegisterCount),
		i:      0x0,
		pc:     0x200,
		sp:     0x0,
		memory: make([]uint16, 4096),
		stack:  make([]uint16, StackSize),
	}

	font := []uint16{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	for i, data := range font {
		cpu.memory[i] = data
	}

	return &cpu
}
