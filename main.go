// CHIP-8 Playground
package main

import (
	"fmt"
	"os"
)

func main() {
	// registers
	// var v = make([]byte, 16)
	var i uint        // 0x000 - 0xFFF
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

	// load program
	memory[pc] = 0xA2
	memory[pc+1] = 0xF0

	for {
		fmt.Printf("I: 0x%X\t\tPC: 0x%X\n", i, pc)

		opscode := uint(memory[pc])<<8 | uint(memory[pc+1])
		opsval := opscode & 0x0FFF

		switch opscode & 0xF000 {
		case 0xA000:
			i = opsval
		}

		pc += 2

		if pc > 0xFFE {
			os.Exit(0)
		}
	}
}
