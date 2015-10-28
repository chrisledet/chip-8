package main

func main() {
	cpu := NewCPU()
	cpu.Clear()

	// load program
	program := []byte{
		0xA2, 0xF0, // set I to 0x2F0
		0x61, 0x67, // set V1 to 0x67
		0xF1, 0x0A, // Wait then store for input at V1
		0x72, 0x32, // add 0x32 to V2
		0x62, 0x35, // set v2 to 0x35
		0x83, 0x10, // copy value from v1 to v3
	}

	cpu.Load(program)
}
