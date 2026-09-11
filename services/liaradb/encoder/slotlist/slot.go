package slotlist

type Slot struct {
	offset int16
	size   int16
}

func NewSlot(offset, size int16) Slot {
	return Slot{offset, size}
}

func (s Slot) Offset() int16 { return s.offset }
func (s Slot) Size() int16   { return s.size }

func (s Slot) Range() (int16, int16) {
	return s.offset, s.offset + s.size
}

func (s Slot) Slice(data []byte) ([]byte, bool) {
	if s.size == 0 {
		return data[s.offset:s.offset], true
	}

	if int(s.offset) >= len(data) {
		return nil, false
	}

	end := s.offset + s.size
	if int(end) > len(data) {
		return nil, false
	}

	return data[s.offset:end], true
}
