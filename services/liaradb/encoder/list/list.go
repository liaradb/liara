package list

import (
	"github.com/liaradb/liaradb/encoder/int32list"
	"github.com/liaradb/liaradb/encoder/raw"
)

type List struct {
	size int32
	list int32list.Int32List
}

func New(data []byte) List {
	l := int32list.New(data)
	size, _ := l.Get(0)
	return List{
		size: size,
		list: l,
	}
}

func (l *List) Length() int {
	return l.list.Length()
}

func (l *List) Size() int32 {
	return l.size
}

func (l *List) setSize(size int32) {
	if l.list.Set(0, size) {
		l.size = size
	}
}

func (l *List) Item(index int32) (int32, bool) {
	return l.list.Get(index + 1)
}

func (l *List) Pop() (int32, error) {
	size := l.Size()
	if size < 1 {
		return 0, raw.ErrUnderflow
	}

	v, ok := l.list.Get(size)
	if !ok {
		return 0, raw.ErrUnderflow
	}

	l.setSize(size - 1)
	return v, nil
}

func (l *List) Push(value int32) (int32, error) {
	size := l.Size()
	if !l.list.Set(size+1, value) {
		return 0, raw.ErrInsufficientSpace
	}

	l.setSize(size + 1)
	return size, nil
}
