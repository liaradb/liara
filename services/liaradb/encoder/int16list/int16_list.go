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

func (l Int16List) Get(index int16) int16 {
	if index >= l.Size() {
		panic("invalid index")
	}

	return int16(binary.BigEndian.Uint16(l.data[l.offset(index):]))
}

func (l Int16List) GetInt32(index int16) int32 {
	if index >= l.Size()-1 {
		panic("invalid index")
	}

	return int32(binary.BigEndian.Uint32(l.data[l.offset(index):]))
}

func (l Int16List) GetInt64(index int16) int64 {
	if index >= l.Size()-3 {
		panic("invalid index")
	}

	return int64(binary.BigEndian.Uint64(l.data[l.offset(index):]))
}

func (l Int16List) Set(index int16, value int16) {
	if index >= l.Size() {
		panic("invalid index")
	}

	binary.BigEndian.PutUint16(l.data[l.offset(index):], uint16(value))
}

func (l Int16List) SetInt32(index int16, value int32) {
	if index >= l.Size()-1 {
		panic("invalid index")
	}

	binary.BigEndian.PutUint32(l.data[l.offset(index):], uint32(value))
}

func (l Int16List) SetInt64(index int16, value int64) {
	if index >= l.Size()-3 {
		panic("invalid index")
	}

	binary.BigEndian.PutUint64(l.data[l.offset(index):], uint64(value))
}

func (l Int16List) offset(index int16) int16 {
	return index * itemSize
}

func (l Int16List) Shift(index, count int16) {
	if index < 0 || count < 0 {
		panic("invalid index")
	}

	if count != 0 {
		copy(l.data[index*itemSize:], l.data[(index-count)*itemSize:])
	}
}

func (l Int16List) ShiftRange(start, end, shift int16) {
	length := end - start

	if start < 0 || length < 0 {
		panic("invalid range")
	}

	if length == 0 || shift == 0 {
		return
	}

	if shift < 0 {
		l.shiftRangeLeft(start, end, -shift)
	} else {
		l.shiftRangeRight(start, end, shift)
	}
}

func (l Int16List) shiftRangeLeft(start, end, shift int16) {
	dstStart := (start - shift) * itemSize
	if dstStart < 0 {
		panic("invalid range")
	}

	container := l.data[dstStart : end*itemSize]
	copy(container, container[shift*itemSize:])
}

func (l Int16List) shiftRangeRight(start, end, shift int16) {
	dstEnd := (end + shift) * itemSize
	if dstEnd > int16(len(l.data)) {
		panic("invalid range")
	}

	container := l.data[start*itemSize : dstEnd]
	copy(container[shift*itemSize:], container)
}
