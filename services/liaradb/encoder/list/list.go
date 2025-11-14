package list

import (
	"github.com/liaradb/liaradb/encoder/int32list"
)

const (
	headerSize = 2
)

type List struct {
	size int32
	next int32
	list int32list.Int32List
}

func New(data []byte) List {
	l := int32list.New(data)
	size, _ := l.Get(0)
	next, _ := l.Get(1)
	return List{
		size: size,
		next: next,
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

func (l *List) Next() int32 {
	return l.next
}

func (l *List) SetNext(next int32) {
	if l.list.Set(1, next) {
		l.next = next
	}
}

func (l *List) Item(index int32) (int32, bool) {
	return l.list.Get(index + headerSize)
}

func (l *List) Pop() (int32, bool) {
	size := l.Size()
	if size < 1 {
		return 0, false
	}

	v, ok := l.list.Get(size + (headerSize - 1))
	if !ok {
		return 0, false
	}

	l.setSize(size - 1)
	return v, true
}

func (l *List) Push(value int32) (int32, bool) {
	size := l.Size()
	if !l.list.Set(size+headerSize, value) {
		return 0, false
	}

	l.setSize(size + 1)
	return size, true
}
