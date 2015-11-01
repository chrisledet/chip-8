package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/chrisledet/chip-8/gfx"
	"github.com/chrisledet/chip-8/io"
)

const (
	StartAddress   int    = 0x200
	MaxProgramSize int    = 0xFFF - StartAddress
	RegisterCount  byte   = 0x10
	StackSize      byte   = 0x10
	FontAddress    uint16 = 0x50
	Width          int    = 64
	Height         int    = 32
)

type CPU struct {
	v  []byte
	i  uint16
	pc uint16

	display [][]bool

	// 0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	// 0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	// 0x200-0xFFF - Program ROM and work RAM
	memory []uint16
	stack  []uint16
	sp     byte

	delayTimer uint16
	soundTimer uint16

	window gfx.Window

	debug bool
}

func (cpu *CPU) Load(program []byte) error {
	if len(program) > MaxProgramSize {
		return errors.New("Program is too larged.")
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

		event := io.Poll()
		if event.IsQuitEvent() {
			break
		}
	}
}

func (cpu *CPU) execute() error {
	opscode := cpu.memory[cpu.pc]<<8 | cpu.memory[cpu.pc+1]
	opsval := opscode & 0x0FFF
	pc := cpu.pc

	if cpu.debug {
		printOpscode(cpu)
	}

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
			// cpu.display.Clear()
			fmt.Printf("TODO - Handle CLEAR...\n")
		}
	case 0xA000:
		cpu.i = opsval
	case 0x1000:
		cpu.pc = opsval - 2
	case 0x2000: // Calls subroutine at NNN.
		// TODO: raise error if cpu.sp is 0xF
		cpu.stack[cpu.sp] = cpu.pc
		cpu.sp++
		cpu.pc = opsval - 2 // reduce by 2 because we increment PC anyway
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
		nn := cpu.memory[cpu.pc+1]
		cpu.v[x] = byte(nn)
	case 0x7000:
		x := cpu.memory[cpu.pc] & 0x0F
		nn := cpu.memory[cpu.pc+1]
		cpu.v[x] += byte(nn)
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
				result = 0x1
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
			if cpu.v[x] > cpu.v[y] {
				cpu.v[0xF] = 0x1
			} else {
				cpu.v[0xF] = 0x0
			}

			cpu.v[x] = cpu.v[y] - cpu.v[x]
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
	case 0xC000:
		x := cpu.memory[cpu.pc] & 0x0F
		nn := cpu.memory[cpu.pc+1]
		rnd := random(0x00, 0xFF)
		result := byte(rnd) & byte(nn)

		cpu.v[x] = result
	case 0xD000:
		x := cpu.memory[cpu.pc] & 0x0F
		y := cpu.memory[cpu.pc+1] / 0x10
		height := cpu.memory[cpu.pc+1] & 0x0F

		posX := uint16(cpu.v[x])
		posY := uint16(cpu.v[y])

		// Draw a sprite at position VX, VY with N bytes of sprite data starting at the address stored in I
		// Set VF to 01 if any set pixels are changed to unset, and 00 otherwise

		cpu.v[0xF] = 0x0

		// cpu.display[height][width]
		for ny := uint16(0); ny < height; ny++ {
			sprite := cpu.memory[cpu.i+ny]
			pixels := byteToDisplay(sprite)

			for nx := uint16(0); nx < 8; nx++ {
				currentPosX := posX + nx
				currentPosY := posY + ny

				if currentPosX >= uint16(Width) {
					currentPosX -= uint16(Width)
				}

				if currentPosY >= uint16(Height) {
					currentPosY -= uint16(Height)
				}

				fmt.Printf("[%d][%d]\n", currentPosX, currentPosY)
				currentPixel := cpu.display[currentPosX][currentPosY]

				if currentPixel != pixels[nx] {
					cpu.display[currentPosX][currentPosY] = true
				} else {
					cpu.v[0xF] = 0x1
				}
			}
		}

		cpu.window.Draw(cpu.display)

	case 0xE000:
		nn := cpu.memory[cpu.pc+1]

		switch nn {
		case 0xA1:
			// Skips the next instruction if the key stored in VX isn't pressed.
		}
	case 0xF000:
		x := cpu.memory[cpu.pc] & 0x0F
		nn := cpu.memory[cpu.pc+1]

		switch nn {
		case 0x07:
			// ignore for now...
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
		case 0x29:
			cpu.i = FontAddress + uint16(cpu.v[x]*5)
		case 0x65:
			for w := uint16(0x0); w <= x; w++ {
				cpu.v[w] = byte(cpu.memory[cpu.i+w])
			}
		case 0x1E:
			cpu.i += uint16(cpu.v[x])
		}
	}

	return nil
}

func NewCPU() *CPU {
	display := make([][]bool, Width, Width)

	for i := 0; i < Width; i++ {
		display[i] = make([]bool, Height, Height)
	}

	cpu := CPU{
		v:       make([]byte, RegisterCount),
		i:       0x0,
		pc:      0x200,
		sp:      0x0,
		memory:  make([]uint16, 4096),
		stack:   make([]uint16, StackSize),
		display: display,
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
		cpu.memory[FontAddress+uint16(i)] = data
	}

	return &cpu
}
