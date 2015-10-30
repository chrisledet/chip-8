package main

import (
	"github.com/chrisledet/c8vm/gfx"
	"github.com/chrisledet/c8vm/io"
)

func main() {
	window := gfx.NewWindow()
	defer window.Close()

	running := true

	for running {
		event := io.Poll()

		if event.IsQuitEvent() {
			running = false
		}
	}
}
