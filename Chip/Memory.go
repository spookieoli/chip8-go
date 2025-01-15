package Chip

type Memory struct {
	Memory [MEMORY_SIZE]uint8
}

// memorySet sets the value at the given index
func (m *Memory) memorySet(index uint16, value uint8) {
	if index < 0 || index >= MEMORY_SIZE {
		panic("Memory index out of bounds")
	}
	m.Memory[index] = value
}

// memoryGet gets the value at the given index
func (m *Memory) memoryGet(index uint16) uint8 {
	if index < 0 || index >= MEMORY_SIZE {
		panic("Memory index out of bounds")
	}
	return m.Memory[index]
}

// memoryGetRange gets a range of values from memory
func (m *Memory) memoryGetRange(start, end uint16) *[]uint8 {
	r := m.Memory[start : start+end]
	return &r
}
