package ProgramLoader

import (
	"os"
)

type ProgramLoader struct {
	filename string
	Program  []byte
}

// NewProgramLoader creates a new program loader
func NewProgramLoader(filename string) *ProgramLoader {
	pl := &ProgramLoader{filename: filename}
	pl.loadProgram()
	return pl
}

// loadProgram loads the program into memory
func (pl *ProgramLoader) loadProgram() {
	// Load the program
	data, err := os.ReadFile(pl.filename)
	if err != nil {
		panic(err)
	}
	pl.Program = data
}
