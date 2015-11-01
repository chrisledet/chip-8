package gfx

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	Title     string = "CHIP-8 VM"
	Width     int32  = 64
	Height    int32  = 32
	PixelSize int32  = 10
	OffColor  uint32 = 0x000000
	OnColor   uint32 = 0xfffffff
)

func init() {
	sdl.Init(sdl.INIT_EVERYTHING)
}

type Window interface {
	Draw([][]bool)
	Close()
}

type SDLWindow struct {
	window *sdl.Window
}

func NewWindow() Window {
	window, err := sdl.CreateWindow(
		Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int(Width*PixelSize),
		int(Height*PixelSize),
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		panic(err)
	}

	return SDLWindow{window: window}
}

func (w SDLWindow) Close() {
	w.window.Destroy()
}

func (w SDLWindow) Draw(pixelState [][]bool) {
	surface, err := w.window.GetSurface()
	if err != nil {
		fmt.Printf("ERROR with GetSurface(): %s\n", err.Error())
		os.Exit(1)
	}

	var pixels []sdl.Rect

	for x := int32(0); x < Width; x++ {
		for y := int32(0); y < Height; y++ {
			color := OffColor
			pixel := sdl.Rect{x * PixelSize, y * PixelSize, PixelSize, PixelSize}
			pixels = append(pixels, pixel)

			if pixelState[x][y] {
				color = OnColor
			}

			surface.FillRect(&pixel, color)
		}
	}

	w.window.UpdateSurfaceRects(pixels)
}
