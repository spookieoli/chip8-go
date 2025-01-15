package window

import (
	"chip-8/Chip"
	"chip-8/ProgramLoader"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

type Window struct {
	Width    int32
	Height   int32
	Title    string
	renderer *sdl.Renderer
	win      *sdl.Window
	pixel    *sdl.Rect
	Chip     *Chip.Chip
}

// NewWindow creates a new window with the given width, height and title
func NewWindow(Title string) *Window {

	// Initialize SDL
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		// We panic if the window can't be created
		panic(err)
	}

	height := int32(Chip.CHIP8_HEIGHT) * Chip.CHIP8_WINDOW_MULTIPLIER
	width := int32(Chip.CHIP8_WIDTH) * Chip.CHIP8_WINDOW_MULTIPLIER

	// Create the window
	win, err := sdl.CreateWindow(Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	// Create the renderer
	renderer, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	// Create the pixel rectangle
	pixel := &sdl.Rect{W: Chip.CHIP8_PIXEL_SIZE, H: Chip.CHIP8_PIXEL_SIZE}

	// create a programLoader
	if len(os.Args) < 2 {
		panic("No program provided - you must provide a Program to run")
	}
	pl := ProgramLoader.NewProgramLoader(os.Args[1])

	// Create the Chip
	chip := Chip.NewChip(&pl.Program)

	// return the window
	return &Window{width, height, Title, renderer, win, pixel, chip}
}

// Show will build and show the SDL Window
func (w *Window) Show() {
	// Here will be the Main Loop
	for {

		// Listen for events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				keyEvent := event.(*sdl.KeyboardEvent)
				if keyEvent.Type == sdl.KEYDOWN {
					w.Chip.Keyboard.MapKeyDown(int(keyEvent.Keysym.Sym))
				} else if keyEvent.Type == sdl.KEYUP {
					w.Chip.Keyboard.MapKeyUp(int(keyEvent.Keysym.Sym))
				}
			}
		}

		// Clear the window
		err := w.renderer.SetDrawColor(0, 0, 0, 0)
		if err != nil {
			fmt.Println(err)
		}

		err = w.renderer.Clear()
		if err != nil {
			fmt.Println(err)
		}

		// Set the draw color
		err = w.renderer.SetDrawColor(255, 255, 255, 255)
		if err != nil {
			fmt.Println(err)
		}

		// Here will be the game logic
		w.logic()

		// Present the window
		w.renderer.Present()

		// Delay if Delaytimer is set
		if w.Chip.Registers.RegistersGetDelayTimer() > 0 {
			sdl.Delay(Chip.CHIP8_DELAY_TIME)
			// substraction of the delay timer
			w.Chip.Registers.RegistersSetDelayTimer(w.Chip.Registers.RegistersGetDelayTimer() - 1)
		}

		// Play the sound if the sound timer is set
		if w.Chip.Registers.RegistersGetSoundTimer() > 0 {
			// Play the sound
			fmt.Print("\a")
			// substraction of the sound timer
			w.Chip.Registers.RegistersSetSoundTimer(w.Chip.Registers.RegistersGetSoundTimer() - 1)
		}

		// Read short from memory
		opcode := w.Chip.ReadShort(w.Chip.Registers.RegistersGetPC())

		// increment the program counter
		w.Chip.Registers.RegisterIncrementPC()

		// Execute the opcode
		w.Chip.ExecuteOpcode(opcode)
	}
}

// logic is the game logic
func (w *Window) logic() {
	// Lets iterate over the screen and draw the pixels
	for i := range w.Chip.Screen.Pixel {
		for j := range w.Chip.Screen.Pixel[i] {
			if w.Chip.Screen.Pixel[i][j] {
				w.pixel.X = int32(i) * Chip.CHIP8_WINDOW_MULTIPLIER
				w.pixel.Y = int32(j) * Chip.CHIP8_WINDOW_MULTIPLIER
				// Draw the pixel
				err := w.renderer.FillRect(w.pixel)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
