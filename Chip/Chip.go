package Chip

import (
	"math/rand/v2"
)

// Constants
const (
	MEMORY_SIZE                         = 4096
	CHIP8_WIDTH                         = 64
	CHIP8_HEIGHT                        = 32
	CHIP8_WINDOW_MULTIPLIER             = 10 // scales the window up
	CHIP8_DATA_REGISTER                 = 16
	CHIP8_STACK_SIZE                    = 16
	CHIP8_TOTAL_KEYS                    = 16
	CHIP8_DEFAULT_CHARACTER_START_POINT = 0x00
	CHIP8_PIXEL_SIZE                    = CHIP8_WINDOW_MULTIPLIER
	CHIP8_DELAY_TIME                    = 1
	CHIP8_DEFAULT_PROGRAM_START_POINT   = 0x200
)

// Chip is the chip itself
type Chip struct {
	Memory    *Memory
	Registers *Registers
	Stack     *Stack
	Keyboard  *Keyboard
	Screen    *Screen
}

// loadDefaultCharacterData loads the default character data into memory
func (c *Chip) loadDefaultCharacterData() {
	// Load the default character data into memory
	characterData := []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	// Load it into memory
	for i, data := range characterData {
		c.Memory.memorySet(uint16(i+CHIP8_DEFAULT_CHARACTER_START_POINT), data)
	}
}

// LoadProgram loads a program into memory
func (c *Chip) LoadProgram(data *[]byte) {

	// Check if the program is too large
	if len(*data) > MEMORY_SIZE-CHIP8_DEFAULT_PROGRAM_START_POINT {
		panic("Program is too large for memory")
	}

	// Load the program into memory
	copy(c.Memory.Memory[CHIP8_DEFAULT_PROGRAM_START_POINT:], *data)

	// set the program counter to the start of the program
	c.Registers.RegistersSetPC(CHIP8_DEFAULT_PROGRAM_START_POINT)
}

// ReadShort reads a short from memory
func (c *Chip) ReadShort(index uint16) uint16 {
	// Read the short from memory and return it - we need to shift the first byte by 8 bits to the left and then OR it with the second byte
	return uint16(c.Memory.memoryGet(index))<<8 | uint16(c.Memory.memoryGet(index+1))
}

// NewChip returns a Chip
func NewChip(program *[]byte) *Chip {
	mem := &Memory{}
	reg := &Registers{}
	stack := &Stack{registers: reg}
	chip := &Chip{mem, reg, stack, NewKeyboard(), NewScreen()}
	chip.loadDefaultCharacterData()
	chip.LoadProgram(program)
	return chip
}

// ExecuteOpcode executes an opcode
func (c *Chip) ExecuteOpcode(opcode uint16) {
	switch opcode {
	case 0x00E0:
		// Clear the screen
		c.Screen.ClearScreen()
	case 0x00EE:
		// Return from a subroutine
		c.Registers.RegistersSetPC(c.Stack.StackPop())
	default:
		c.ExecuteExtended(opcode)
	}
}

// ExecuteExtended executes an extended opcode
func (c *Chip) ExecuteExtended(opcode uint16) {
	// ignore the first 4 bits
	nnn := opcode & 0x0FFF
	// take the last 4 bits
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	kk := opcode & 0x00FF
	// here we will be using the first 4 bits
	switch opcode & 0xF000 {
	case 0x1000:
		// Jump to location nnn
		c.Registers.RegistersSetPC(nnn)
	case 0x2000:
		// Call subroutine at nnn
		c.Stack.StackPush(c.Registers.RegistersGetPC())
		c.Registers.RegistersSetPC(nnn)
	case 0x3000:
		// Skip next instruction if Vx = kk
		kk := opcode & 0x00FF
		if uint16(c.Registers.RegistersGetV(x)) == kk {
			c.Registers.RegisterIncrementPC()
		}
	case 0x4000:
		// Skip next instruction if Vx != kk
		if uint16(c.Registers.RegistersGetV(x)) != kk {
			c.Registers.RegisterIncrementPC()
		}
	case 0x5000:
		// Skip next instruction if Vx = Vy
		if c.Registers.RegistersGetV(x) == c.Registers.RegistersGetV(y) {
			c.Registers.RegisterIncrementPC()
		}
	case 0x6000:
		// Set Vx = kk
		c.Registers.RegistersSetV(int(x), uint8(kk))
	case 0x7000:
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(x)+uint8(kk))
	case 0x8000:
		c.ExecutedExtendedExtra(opcode)
	case 0x9000:
		// Skip next instruction if Vx != Vy
		if c.Registers.RegistersGetV(x) != c.Registers.RegistersGetV(y) {
			c.Registers.RegisterIncrementPC()
		}
	case 0xA000:
		// Set I = nnn
		c.Registers.RegistersSetI(nnn)
	case 0xB000:
		// Jump to location nnn + V0
		c.Registers.RegistersSetPC(nnn + uint16(c.Registers.RegistersGetV(0)))
	case 0xC000:
		// Set Vx = random byte AND kk
		c.Registers.RegistersSetV(int(x), uint8(rand.UintN(256))&uint8(kk))
	case 0xD000:
		// get n
		n := opcode & 0x000F
		// Draw a sprite at position Vx, Vy with n bytes of sprite data starting at the address stored in I
		collision := c.Screen.DrawSprite(int(c.Registers.RegistersGetV(x)), int(c.Registers.RegistersGetV(y)), c.Memory.memoryGetRange(c.Registers.RegistersGetI(), n))
		if collision {
			c.Registers.RegistersSetV(0xF, 1)
		} else {
			c.Registers.RegistersSetV(0xF, 0)
		}
	case 0xE000:
		c.ExecuteOpcodeExtendedE(opcode)
	case 0xF000:
		c.ExecuteOpcodeExtendedF(opcode)
	}
}

// ExecuteOpcodeExtendedF This will execute the extended opcode 0xF000
func (c *Chip) ExecuteOpcodeExtendedF(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	switch opcode & 0x00FF {
	case 0x0007:
		// Set Vx = delay timer value
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetDelayTimer())
	case 0x000A:
		// Wait for a key press, store the value of the key in Vx
		pressed := false
		for i := 0; i < CHIP8_TOTAL_KEYS; i++ {
			if c.Keyboard.IsKeyDown(i) {
				c.Registers.RegistersSetV(int(x), uint8(i))
				pressed = true
			}
		}
		if !pressed {
			c.Registers.RegisterDecrementPC()
		}
	case 0x0015:
		// Set delay timer = Vx
		c.Registers.RegistersSetDelayTimer(c.Registers.RegistersGetV(x))
	case 0x0018:
		// Set sound timer = Vx
		c.Registers.RegistersSetSoundTimer(c.Registers.RegistersGetV(x))
	case 0x001E:
		// Set I = I + Vx
		c.Registers.RegistersSetI(c.Registers.RegistersGetI() + uint16(c.Registers.RegistersGetV(x)))
	case 0x0029:
		// Set I = location of sprite for digit Vx
		c.Registers.RegistersSetI(uint16(c.Registers.RegistersGetV(x)) * 5)
	case 0x0033:
		// Store BCD representation of Vx in memory locations I, I+1, and I+2
		vx := c.Registers.RegistersGetV(x)
		c.Memory.memorySet(c.Registers.RegistersGetI(), vx/100)
		c.Memory.memorySet(c.Registers.RegistersGetI()+1, (vx/10)%10)
		c.Memory.memorySet(c.Registers.RegistersGetI()+2, vx%10)
	case 0x0055:
		// Store registers V0 through Vx in memory starting at location I
		for i := 0; i <= int(x); i++ {
			c.Memory.memorySet(c.Registers.RegistersGetI()+uint16(i), c.Registers.RegistersGetV(uint16(i)))
		}
	case 0x0065:
		// Read registers V0 through Vx from memory starting at location I
		for i := 0; i <= int(x); i++ {
			c.Registers.RegistersSetV(i, c.Memory.memoryGet(c.Registers.RegistersGetI()+uint16(i)))
		}
	}
}

// ExecuteOpcodeExtendedE This will execute the extended opcode 0xE000
func (c *Chip) ExecuteOpcodeExtendedE(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	switch opcode & 0x00FF {
	case 0x009E:
		// Skip next instruction if key with the value of Vx is pressed
		if c.Keyboard.IsKeyDown(int(c.Registers.RegistersGetV(x))) {
			c.Registers.RegisterIncrementPC()
		}
	case 0x00A1:
		// Skip next instruction if key with the value of Vx is not pressed
		if !c.Keyboard.IsKeyDown(int(c.Registers.RegistersGetV(x))) {
			c.Registers.RegisterIncrementPC()
		}
	}
}

// ExecutedExtendedExtra executes an extended extra opcode 0x8000
func (c *Chip) ExecutedExtendedExtra(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	switch opcode & 0x000F {
	case 0x0000:
		// Set Vx = Vy
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(y))
	case 0x0001:
		// Set Vx = Vx OR Vy
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(x)|c.Registers.RegistersGetV(y))
	case 0x0002:
		// Set Vx = Vx AND Vy
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(x)&c.Registers.RegistersGetV(y))
	case 0x0003:
		// Set Vx = Vx XOR Vy
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(x)^c.Registers.RegistersGetV(y))
	case 0x0004:
		// Add Vx, Vy
		sum := uint16(c.Registers.RegistersGetV(x)) + uint16(c.Registers.RegistersGetV(y))
		if sum > 0xFF {
			c.Registers.RegistersSetV(0xF, 1)
		} else {
			c.Registers.RegistersSetV(0xF, 0)
		}
		// store only the lowest 8 bits
		c.Registers.RegistersSetV(int(x), uint8(sum))

	case 0x0005:
		// Subtract Vy from Vx
		if c.Registers.RegistersGetV(x) > c.Registers.RegistersGetV(y) {
			c.Registers.RegistersSetV(0xF, 1)
		} else {
			c.Registers.RegistersSetV(0xF, 0)
		}
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(x)-c.Registers.RegistersGetV(y))
	case 0x0006:
		// Store the least significant bit of Vx in VF and then shift Vx to the right by 1
		c.Registers.RegistersSetV(0xF, c.Registers.RegistersGetV(x)&0x1)
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(x)>>1)
	case 0x0007:
		// Set Vx = Vy - Vx
		if c.Registers.RegistersGetV(y) > c.Registers.RegistersGetV(x) {
			c.Registers.RegistersSetV(0xF, 1)
		} else {
			c.Registers.RegistersSetV(0xF, 0)
		}
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(y)-c.Registers.RegistersGetV(x))
	case 0x000E:
		// Store the most significant bit of Vx in VF and then shift Vx to the left by 1
		c.Registers.RegistersSetV(0xF, c.Registers.RegistersGetV(x)>>7)
		c.Registers.RegistersSetV(int(x), c.Registers.RegistersGetV(x)<<1)
	}
}
