// CHIP-8 Playground
package main

import (
	"fmt"
	"os"
)

func main() {
	// registers
	var v = make([]byte, 16)
	var i uint = 0x0  // 0x000 - 0xFFF
	var pc uint = 0x0 // 0x000 - 0xFFF 2 bytes

	// timers
	// var delayTimer int
	// var soundTimer int

	// storage
	var memory = make([]byte, 4096)
	// var key = make([]uint, 16)
	// var stack = make([]uint, 16)
	// var sp uint

	// 0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	// 0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	// 0x200-0xFFF - Program ROM and work RAM
	// var gfx = make([]byte, 2048)

	// reset timers
	pc = 0x200

	// clear
	for x, _ := range v {
		v[x] = 0
	}

	// load program
	memory[pc] = 0xA2
	memory[pc+1] = 0xF0
	memory[pc+2] = 0xA2
	memory[pc+3] = 0xF5
	memory[pc+4] = 0x61
	memory[pc+5] = 0x4F
	memory[pc+6] = 0x62
	memory[pc+7] = 0x35
	memory[pc+8] = 0x6F
	memory[pc+9] = 0x35

	for {
		opscode := uint(memory[pc])<<8 | uint(memory[pc+1])
		opsval := opscode & 0x0FFF

		switch opscode & 0xF000 {
		case 0xA000:
			previous := i
			i = opsval

			fmt.Printf("Setting I from 0x%X to 0x%X - %x\n", previous, opsval, i)
		case 0x1000:
			pc = opsval
			continue
		case 0x6000:
			x := uint(memory[pc]) & 0x0F
			opsval = uint(memory[pc+1])
			v[x] = byte(opsval)

			fmt.Printf("Setting V%d to 0x%X\n", x, opsval)
		}

		pc += 2

		if pc > 0xFFE {
			os.Exit(0)
		}
	}
}
