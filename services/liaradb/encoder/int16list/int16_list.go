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

// TODO: Test this
func (l Int16List) GetInt32(index int16) (int32, bool) {
	if index >= l.Size()+1 {
		return 0, false
	}

	return int32(binary.BigEndian.Uint32(l.data[l.offset(index):])), true
}

// TODO: Test this
func (l Int16List) GetInt64(index int16) (int64, bool) {
	if index >= l.Size()+3 {
		return 0, false
	}

	return int64(binary.BigEndian.Uint64(l.data[l.offset(index):])), true
}

func (l Int16List) Set(index int16, value int16) bool {
	if index >= l.Size() {
		return false
	}

	binary.BigEndian.PutUint16(l.data[l.offset(index):], uint16(value))
	return true
}

// TODO: Test this
func (l Int16List) Set32(index int16, value int32) bool {
	if index >= l.Size()+1 {
		return false
	}

	binary.BigEndian.PutUint32(l.data[l.offset(index):], uint32(value))
	return true
}

// TODO: Test this
func (l Int16List) Set64(index int16, value int64) bool {
	if index >= l.Size()+3 {
		return false
	}

	binary.BigEndian.PutUint64(l.data[l.offset(index):], uint64(value))
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

func (l Int16List) ShiftRange(start, end, shift int16) bool {
	length := end - start

	if start < 0 || length < 0 {
		return false
	}

	if length == 0 || shift == 0 {
		return true
	}

	if shift < 0 {
		return l.shiftRangeLeft(start, end, -shift)
	} else {
		return l.shiftRangeRight(start, end, shift)
	}
}

func (l Int16List) shiftRangeLeft(start, end, shift int16) bool {
	dstStart := (start - shift) * itemSize
	if dstStart < 0 {
		return false
	}

	container := l.data[dstStart : end*itemSize]
	copy(container, container[shift*itemSize:])

	return true
}

func (l Int16List) shiftRangeRight(start, end, shift int16) bool {
	dstEnd := (end + shift) * itemSize
	if dstEnd > int16(len(l.data)) {
		return false
	}

	container := l.data[start*itemSize : dstEnd]
	copy(container[shift*itemSize:], container)

	return true
}
