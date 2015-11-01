package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chrisledet/chip-8/gfx"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Must specify ROM file")
		os.Exit(1)
	}

	romPath := os.Args[1]
	program, err := ioutil.ReadFile(romPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem with loading rom: %s\n", err.Error())
		os.Exit(1)
	}

	cpu := NewCPU()
	cpu.debug = true

	gfxWindow := gfx.NewWindow()
	cpu.window = gfxWindow

	cpu.Load(program)

	gfxWindow.Close()
}
