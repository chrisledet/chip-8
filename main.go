// CHIP-8 Playground
package main

import (
	"fmt"
	"os"
)

func main() {
	// registers
	var v = make([]byte, 16)
	var i uint16 = 0x0  // 0x000 - 0xFFF
	var pc uint16 = 0x0 // 0x000 - 0xFFF 2 bytes

	// timers
	// var delayTimer int
	// var soundTimer int

	// storage
	var memory = make([]uint16, 4096)
	// var key = make([]uint16, 16)
	var stack = make([]uint16, 16)
	var sp byte

	// 0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	// 0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	// 0x200-0xFFF - Program ROM and work RAM
	// var gfx = make([]byte, 2048)

	// reset timers
	pc = 0x200

	// clear
	sp = 0
	for x, _ := range v {
		v[x] = 0
	}

	// load program
	memory[pc] = 0xA2
	memory[pc+1] = 0xF0

	memory[pc+2] = 0xA2
	memory[pc+3] = 0xF5

	memory[pc+4] = 0x61
	memory[pc+5] = 0x67

	memory[pc+6] = 0x62
	memory[pc+7] = 0x35

	memory[pc+8] = 0x72
	memory[pc+9] = 0x32

	memory[pc+10] = 0x22
	memory[pc+11] = 0x0C

	memory[pc+12] = 0x00
	memory[pc+13] = 0xEE

	memory[pc+14] = 0x51
	memory[pc+15] = 0x20

	memory[pc+16] = 0x62
	memory[pc+17] = 0x35

	memory[pc+18] = 0x83
	memory[pc+19] = 0x10

	for {
		opscode := memory[pc]<<8 | memory[pc+1]
		opsval := opscode & 0x0FFF

		if opscode == 0x00EE {
			if sp > 0 {
				pc = stack[sp-1]
				sp--
				fmt.Printf("CMD: subroutine return to address 0x%X\n", pc)
			} else {
				fmt.Printf("ERROR: invalid subroutine return at address 0x%X\n", pc)
			}
		}

		switch opscode & 0xF000 {
		case 0xA000:
			i = opsval
			fmt.Printf("CMD: set register I to 0x%X\n", i)
		case 0x1000:
			pc = opsval
			continue
		case 0x2000: // Calls subroutine at NNN.
			stack[sp] = pc
			sp++
			pc = opsval
			fmt.Printf("CMD: jump to subroutine at address 0x%X\n", opsval)
			continue
		case 0x3000:
			x := memory[pc] & 0x0F
			nn := memory[pc+1]

			fmt.Printf("CMD: skip next instruction if V%d (0x%X) = 0x%X\n", x, v[x], nn)
			if v[x] == byte(nn) {
				pc += 2
			}
		case 0x4000:
			x := memory[pc] & 0x0F
			nn := memory[pc+1]

			fmt.Printf("CMD: skip next instruction if V%d (0x%X) != 0x%X\n", x, v[x], nn)
			if v[x] != byte(nn) {
				pc += 2
			}
		case 0x5000:
			x := memory[pc] & 0x0F
			y := memory[pc+1] / 0x10

			fmt.Printf("CMD: skip next instruction if V%d (0x%X) != V%d (0x%X)\n", x, v[x], y, v[y])
			if v[x] == v[y] {
				pc += 2
			}
		case 0x6000:
			x := memory[pc] & 0x0F
			opsval = memory[pc+1]
			v[x] = byte(opsval)

			fmt.Printf("CMD: set V%d to 0x%X\n", x, opsval)
		case 0x7000:
			x := memory[pc] & 0x0F
			opsval = memory[pc+1]
			v[x] += byte(opsval)

			fmt.Printf("CMD: add 0x%X to V%d\n", opsval, x)
		case 0x8000:
			x := memory[pc] & 0x0F
			y := memory[pc+1] / 0x10

			fmt.Printf("CMD: set V%d (0x%X) to V%d (0x%X)\n", x, v[x], y, v[y])
			v[x] = v[y]
		}

		pc += 2

		if pc > 0xFFE {
			os.Exit(0)
		}
	}
}
