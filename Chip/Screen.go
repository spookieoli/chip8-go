package Chip

// Screen Struct representing the screen
type Screen struct {
	Pixel [CHIP8_WIDTH][CHIP8_HEIGHT]bool
}

// NewScreen creates a new screen
func NewScreen() *Screen {
	return &Screen{}
}

// SetPixel sets the pixel at the given x and y coordinates
func (s *Screen) SetPixel(x, y int) {
	s.Pixel[x][y] = !s.Pixel[x][y]
}

// IsPixelSet returns true if the pixel at the given x and y coordinates is set
func (s *Screen) IsPixelSet(x, y int) bool {
	return s.Pixel[x][y]
}

// DrawSprite draws an 8-bit wide sprite at the specified (x, y) coordinates and returns true if any pixels were erased.
func (s *Screen) DrawSprite(x, y int, sprite *[]uint8) (collision bool) {
	for idx, v := range *sprite {
		for j := 0; j < 8; j++ {
			if v&(0x80>>j) == 0 {
				continue // Skip this pixel
			}

			// Check for collision
			if s.Pixel[(j+x)%CHIP8_WIDTH][(idx+y)%CHIP8_HEIGHT] {
				collision = true
			}

			// Calculate the x and y coordinates
			s.Pixel[(j+x)%CHIP8_WIDTH][(idx+y)%CHIP8_HEIGHT] = !s.Pixel[(j+x)%CHIP8_WIDTH][(idx+y)%CHIP8_HEIGHT]
		}
	}
	return
}

// ClearScreen clears the screen
func (s *Screen) ClearScreen() {
	// Clear the screen by copying an empty array into the pixel array
	copy(s.Pixel[:], make([][CHIP8_HEIGHT]bool, CHIP8_WIDTH))
}
