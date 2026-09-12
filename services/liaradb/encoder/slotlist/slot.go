package slotlist

type Slot struct {
	offset int16
	size   int16
	data   []byte
}

// TODO: Make this private
func NewSlot(offset, size int16, data []byte) Slot {
	return Slot{offset, size, data}
}

func (s Slot) Offset() int16 { return s.offset }
func (s Slot) Size() int16   { return s.size }

func (s Slot) Slice() ([]byte, bool) {
	if s.size == 0 {
		return s.data[s.offset:s.offset], true
	}

	if int(s.offset) >= len(s.data) {
		return nil, false
	}

	end := s.offset + s.size
	if int(end) > len(s.data) {
		return nil, false
	}

	return s.data[s.offset:end], true
}

func (s Slot) SliceUnsafe() []byte {
	return s.data[s.offset : s.offset+s.size]
}
