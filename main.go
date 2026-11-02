package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rivo/tview"

	"github.com/cjr29/zen80/io"
	"github.com/cjr29/zen80/memory"
	"github.com/cjr29/zen80/z80"
)

func main() {
	// Display gridview
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	menu := newPrimitive("Menu")
	main := newPrimitive("Main content")
	sideBar := newPrimitive("Side Bar")
	footer := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(footer, 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false).
		AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	if err := tview.NewApplication().SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

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
	fmt.Fprintf(footer, "Result: %d\n", cpu.A)
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
	fmt.Printf("Loaded program %q (%d bytes)\n", path, len(program))
	return program, nil
}



