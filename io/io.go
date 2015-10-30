package io

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Event interface {
	IsQuitEvent() bool
}

type SDLEvent struct {
	event sdl.Event
}

func (e SDLEvent) IsQuitEvent() bool {
	sdlEvent := e.event

	switch sdlEvent.(type) {
	case *sdl.QuitEvent:
		return true
	default:
		return false
	}
}

func Poll() *SDLEvent {
	event := sdl.PollEvent()
	return &SDLEvent{event: event}
}
