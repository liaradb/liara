package list

import (
	"encoding/binary"

	"github.com/liaradb/liaradb/encoder/raw"
)

const (
	headerSize = 4
	itemSize   = 4
)

type List struct {
	size int32
	data []byte
}

func New(data []byte) *List {
	return &List{
		size: getInt32(data, 0),
		data: data,
	}
}

func getInt32(data []byte, offset int32) int32 {
	if len(data) < int(offset+4) {
		return 0
	}

	return int32(binary.BigEndian.Uint32(data[offset:]))
}

func setInt32(data []byte, offset int32, value int32) bool {
	if len(data) < int(offset+4) {
		return false
	}

	binary.BigEndian.PutUint32(data[offset:], uint32(value))
	return true
}

func (l *List) Length() int {
	return len(l.data)
}

func (l *List) Size() int32 {
	return l.size
}

func (l *List) setSize(size int32) {
	if setInt32(l.data, 0, size) {
		l.size = size
	}
}

func (l *List) value(offset int32) (int32, bool) {
	if len(l.data) < int(offset)+4 {
		return 0, false
	}

	return int32(binary.BigEndian.Uint32(l.data[offset:])), true
}

func (l *List) Item(index int32) (int32, bool) {
	return l.value(l.offset(index))
}

func (l *List) setValue(offset int32, value int32) bool {
	return setInt32(l.data, offset, value)
}

func (l *List) Append(value int32) (int, error) {
	size := l.Size()
	offset := l.offset(size)
	if !l.setValue(offset, value) {
		return 0, raw.ErrInsufficientSpace
	}

	l.setSize(size + 1)
	return int(size), nil
}

func (l *List) offset(index int32) int32 {
	return headerSize + index*itemSize
}
