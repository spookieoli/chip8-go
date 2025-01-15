package Chip

type Stack struct {
	stack     [CHIP8_STACK_SIZE]uint16
	registers *Registers
}

// StackPush pushes a value onto the stack
func (s *Stack) StackPush(value uint16) {
	s.stack[s.registers.sp] = value
	s.registers.RegistersIncrementSP()
}

// StackPop pops a value off the stack
func (s *Stack) StackPop() uint16 {
	s.registers.RegistersDecrementSP()
	return s.stack[s.registers.sp]
}
