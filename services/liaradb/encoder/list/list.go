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
	data []byte
}

func New(data []byte) *List {
	return &List{
		data: data,
	}
}

func (l *List) Length() int {
	return len(l.data)
}

func (l *List) Size() int32 {
	if l.Length() < 4 {
		return 0
	}

	return int32(binary.BigEndian.Uint32(l.data))
}

func (l *List) setSize(value int32) {
	binary.BigEndian.PutUint32(l.data, uint32(value))
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
	if len(l.data) < int(offset)+4 {
		return false
	}

	binary.BigEndian.PutUint32(l.data[offset:], uint32(value))

	return true
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
