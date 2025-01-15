package Chip

import (
	sdl "github.com/veandco/go-sdl2/sdl"
)

type Keyboard struct {
	Keys   [CHIP8_TOTAL_KEYS]bool
	KeyMap map[int]int
}

// NewKeyboard creates a new Keyboard
func NewKeyboard() *Keyboard {
	keyMap := map[int]int{
		sdl.K_0: 0x0,
		sdl.K_1: 0x1,
		sdl.K_2: 0x2,
		sdl.K_3: 0x3,
		sdl.K_4: 0x4,
		sdl.K_5: 0x5,
		sdl.K_6: 0x6,
		sdl.K_7: 0x7,
		sdl.K_8: 0x8,
		sdl.K_9: 0x9,
		sdl.K_a: 0xA,
		sdl.K_b: 0xB,
		sdl.K_c: 0xC,
		sdl.K_d: 0xD,
		sdl.K_e: 0xE,
		sdl.K_f: 0xF,
	}
	return &Keyboard{KeyMap: keyMap}
}

// MapKeyDown maps the key down
func (k *Keyboard) MapKeyDown(key int) {
	k.keyDown(k.keyMapper(key))
}

// MapKeyUp maps the key up
func (k *Keyboard) MapKeyUp(key int) {
	k.keyUp(k.keyMapper(key))
}

// KeyMapper maps the keys to the chip8 keys
func (k *Keyboard) keyMapper(key int) int {
	if mappedKey, ok := k.KeyMap[key]; ok {
		return mappedKey
	}
	return -1
}

// KeyDown sets the key to down
func (k *Keyboard) keyDown(key int) {
	if key == -1 {
		return
	}
	k.Keys[key] = true
}

// KeyUp sets the key to up
func (k *Keyboard) keyUp(key int) {
	if key == -1 {
		return
	}
	k.Keys[key] = false
}

// IsKeyDown returns true if the key is down
func (k *Keyboard) IsKeyDown(key int) bool {
	if key == -1 {
		return false
	}
	return k.Keys[key]
}

// IsKeyUp returns true if the key is up
func (k *Keyboard) IsKeyUp(key int) bool {
	if key == -1 {
		return false
	}
	return !k.Keys[key]
}
