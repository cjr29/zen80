package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProgramFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "program.bin")
	want := []byte{0x3E, 0x05, 0x76}
	if err := os.WriteFile(path, want, 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := loadProgramFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("loadProgramFromFile() = %v, want %v", got, want)
	}
}

func TestLoadProgramFromFileRejectsOversizedProgram(t *testing.T) {
	path := filepath.Join(t.TempDir(), "program.bin")
	if err := os.WriteFile(path, make([]byte, 1<<16+1), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := loadProgramFromFile(path); err == nil {
		t.Fatal("loadProgramFromFile() accepted a program larger than 64 KB")
	}
}
