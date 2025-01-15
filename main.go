package main

import (
	"chip-8/window"
	_ "github.com/veandco/go-sdl2/sdl"
)

// System main function
func main() {
	// Create a new window
	win := window.NewWindow("Chip-8 Emulator")
	win.Show()
}
