package main

import (
	"math/rand"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func byteToDisplay(b uint16) [8]bool {
	display := [8]bool{
		b&0x80 == 0x80,
		b&0x40 == 0x40,
		b&0x20 == 0x20,
		b&0x10 == 0x10,
		b&0x08 == 0x08,
		b&0x04 == 0x04,
		b&0x02 == 0x02,
		b&0x01 == 0x01,
	}

	return display
}
