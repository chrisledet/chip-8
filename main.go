package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chrisledet/chip-8/gfx"
)

func main() {
	cpu := NewCPU()
	cpu.debug = true

	gfxWindow := gfx.NewWindow()
	cpu.window = gfxWindow

	program, err := ioutil.ReadFile("roms/MERLIN")
	if err != nil {
		fmt.Errorf("Problem with loading rom: %s\n", err.Error())
		os.Exit(1)
	}

	cpu.Load(program)

	gfxWindow.Close()
}
