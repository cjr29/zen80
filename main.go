package main

import (
	"fmt"

	"github.com/ha1tch/zen80/io"
	"github.com/ha1tch/zen80/memory"
	"github.com/ha1tch/zen80/z80"
)

func main() {
	// Create memory and I/O
	mem := memory.NewRAM()
	io := io.NewNullIO()

	// Load a program
	program := []uint8{
		0x3E, 0x05, // LD A, 5
		0x06, 0x03, // LD B, 3
		0x80, // ADD A, B
		0x76, // HALT
	}
	mem.Load(0x0000, program)

	// Create and run CPU
	cpu := z80.New(mem, io)
	for !cpu.Halted {
		cpu.Step()
	}

	// Result is in register A
	fmt.Printf("Result: %d\n", cpu.A)
}
