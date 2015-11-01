package main

import (
	"testing"
)

func TestByteToDisplay(t *testing.T) {
	var data [8]bool

	data = byteToDisplay(128)
	if data[0] != true {
		t.Errorf("0b10000000 failed")
	}

	data = byteToDisplay(64)
	if data[1] != true {
		t.Errorf("0b01000000 failed")
	}

	data = byteToDisplay(32)
	if data[2] != true {
		t.Errorf("0b00100000 failed")
	}

	data = byteToDisplay(16)
	if data[3] != true {
		t.Errorf("0b00010000 failed")
	}

	data = byteToDisplay(8)
	if data[4] != true {
		t.Errorf("0b00001000 failed")
	}

	data = byteToDisplay(4)
	if data[5] != true {
		t.Errorf("0b00000100 failed")
	}

	data = byteToDisplay(2)
	if data[6] != true {
		t.Errorf("0b00000010 failed")
	}

	data = byteToDisplay(1)
	if data[7] != true {
		t.Errorf("0b00000001 failed")
	}
}
