package Chip

type Registers struct {
	v          [CHIP8_DATA_REGISTER]uint8
	i          uint16
	delayTimer uint8
	soundTimer uint8
	pc         uint16
	sp         uint8
}

// RegistersGetV returns the value of the register at the given index
func (r *Registers) RegistersGetV(index uint16) uint8 {
	if index < 0 || index >= CHIP8_DATA_REGISTER {
		panic("Register index out of bounds")
	}
	return r.v[index]
}

// RegistersSetV sets the value of the register at the given index
func (r *Registers) RegistersSetV(index int, value uint8) {
	if index < 0 || index >= CHIP8_DATA_REGISTER {
		panic("Register index out of bounds")
	}
	r.v[index] = value
}

// RegistersGetI returns the value of the I register
func (r *Registers) RegistersGetI() uint16 {
	return r.i
}

// RegistersSetI sets the value of the I register
func (r *Registers) RegistersSetI(value uint16) {
	r.i = value
}

// RegistersGetDelayTimer returns the value of the delay timer
func (r *Registers) RegistersGetDelayTimer() uint8 {
	return r.delayTimer
}

// RegistersSetDelayTimer sets the value of the delay timer
func (r *Registers) RegistersSetDelayTimer(value uint8) {
	r.delayTimer = value
}

// RegistersGetSoundTimer returns the value of the sound timer
func (r *Registers) RegistersGetSoundTimer() uint8 {
	return r.soundTimer
}

// RegistersSetSoundTimer sets the value of the sound timer
func (r *Registers) RegistersSetSoundTimer(value uint8) {
	r.soundTimer = value
}

// RegistersGetPC returns the value of the program counter
func (r *Registers) RegistersGetPC() uint16 {
	return r.pc
}

// RegistersSetPC sets the value of the program counter
func (r *Registers) RegistersSetPC(value uint16) {
	r.pc = value
}

// RegistersGetSP returns the value of the stack pointer
func (r *Registers) RegistersGetSP() uint8 {
	return r.sp
}

// RegistersSetSP sets the value of the stack pointer
func (r *Registers) RegistersSetSP(value uint8) {
	r.sp = value
}

// RegistersIncrementSP increments the Stackpointer
func (r *Registers) RegistersIncrementSP() {
	r.sp++
	if r.sp > CHIP8_STACK_SIZE {
		panic("Stack overflow")
	}
}

// RegistersDecrementSP decrements the Stackpointer
func (r *Registers) RegistersDecrementSP() {
	r.sp--
}

// RegisterIncrementPC increments the program counter
func (r *Registers) RegisterIncrementPC() {
	r.pc += 2
}

// RegisterDecrementPC decrements the program counter
func (r *Registers) RegisterDecrementPC() {
	r.pc -= 2
}
