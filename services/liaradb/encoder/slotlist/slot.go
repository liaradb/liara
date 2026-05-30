package slotlist

type Slot struct {
	offset int16
	size   int16
}

func (s Slot) Offset() int16 { return s.offset }
func (s Slot) Size() int16   { return s.size }

func (s Slot) Range() (int16, int16) {
	return s.offset, s.offset + s.size
}
