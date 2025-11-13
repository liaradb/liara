package list

import "encoding/binary"

const (
	itemSize = 4
)

type Int32List struct {
	data []byte
}

func New(data []byte) Int32List {
	return Int32List{
		data: data,
	}
}

func (l Int32List) Length() int {
	return len(l.data)
}

func (l Int32List) Size() int32 {
	return int32(len(l.data) / itemSize)
}

func (l Int32List) Get(index int32) (int32, bool) {
	if index >= l.Size() {
		return 0, false
	}

	return int32(binary.BigEndian.Uint32(l.data[l.offset(index):])), true
}

func (l Int32List) Set(index int32, value int32) bool {
	if index >= l.Size() {
		return false
	}

	binary.BigEndian.PutUint32(l.data[l.offset(index):], uint32(value))
	return true
}

func (l Int32List) offset(index int32) int32 {
	return index * itemSize
}
