package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ha1tch/zen80/io"
	"github.com/ha1tch/zen80/memory"
	"github.com/ha1tch/zen80/z80"
)

func main() {
	programPath := flag.String("program", "", "path to a raw Z80 binary to run")
	flag.Parse()

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
	if *programPath != "" {
		var err error
		program, err = loadProgramFromFile(*programPath)
		if err != nil {
			log.Fatal(err)
		}
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

func loadProgramFromFile(path string) ([]uint8, error) {
	program, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read program %q: %w", path, err)
	}
	if len(program) > 1<<16 {
		return nil, fmt.Errorf("program %q is %d bytes; maximum size is 65536 bytes", path, len(program))
	}
	return program, nil
}
