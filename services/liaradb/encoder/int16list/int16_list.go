package int16list

import "encoding/binary"

const (
	itemSize = 2
)

type Int16List struct {
	data []byte
}

func New(data []byte) Int16List {
	return Int16List{
		data: data,
	}
}

func (l Int16List) Length() int {
	return len(l.data)
}

func (l Int16List) Size() int16 {
	return int16(len(l.data) / itemSize)
}

func (l Int16List) Get(index int16) (int16, bool) {
	if index >= l.Size() {
		return 0, false
	}

	return int16(binary.BigEndian.Uint16(l.data[l.offset(index):])), true
}

func (l Int16List) Set(index int16, value int16) bool {
	if index >= l.Size() {
		return false
	}

	binary.BigEndian.PutUint16(l.data[l.offset(index):], uint16(value))
	return true
}

func (l Int16List) offset(index int16) int16 {
	return index * itemSize
}

func (l Int16List) Shift(index, count int16) bool {
	if index < 0 || count < 0 {
		return false
	}

	if count == 0 {
		return true
	}

	copy(l.data[index*itemSize:], l.data[(index-count)*itemSize:])
	return true
}
