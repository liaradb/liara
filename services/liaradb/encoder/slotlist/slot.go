package slotlist

type Slot struct {
	offset int16
	size   int16
}

func (s Slot) Offset() int16  { return s.offset }
func (s Slot) Size() int16    { return s.size }
func (s Slot) IsFilled() bool { return s.offset > 0 }

func (s Slot) Range(headerSize int16) (int16, int16) {
	start := headerSize + s.offset
	end := start + s.size
	return start, end
}
