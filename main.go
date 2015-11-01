package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chrisledet/chip-8/gfx"
)

func main() {
	romPath := os.Args[1]
	program, err := ioutil.ReadFile(romPath)
	if err != nil {
		fmt.Errorf("Problem with loading rom: %s\n", err.Error())
		os.Exit(1)
	}

	cpu := NewCPU()
	cpu.debug = true

	gfxWindow := gfx.NewWindow()
	cpu.window = gfxWindow

	cpu.Load(program)

	gfxWindow.Close()
}
