package main

import (
	"fmt"
)

func printOpscode(cpu *CPU) {
	fmt.Printf("EXECUTING: ")

	if cpu.memory[cpu.pc] < 0x10 {
		fmt.Printf("0x0%X", cpu.memory[cpu.pc])

		if cpu.memory[cpu.pc+1] < 0x10 {
			fmt.Printf("0%X", cpu.memory[cpu.pc+1])
		} else {
			fmt.Printf("%X", cpu.memory[cpu.pc+1])
		}
	} else {
		fmt.Printf("0x%X", cpu.memory[cpu.pc])

		if cpu.memory[cpu.pc+1] < 0x10 {
			fmt.Printf("0%X", cpu.memory[cpu.pc+1])
		} else {
			fmt.Printf("%X", cpu.memory[cpu.pc+1])
		}
	}

	fmt.Printf(", ADDR: 0x%X\n", cpu.pc)
}
