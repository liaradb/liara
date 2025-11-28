package list

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/int32list"
)

const (
	headerSize = 2
	itemSize   = 4
)

type List struct {
	size int32
	next int32
	list int32list.Int32List
}

func New(data []byte) List {
	l := int32list.New(data)
	size, _ := l.Get(0)

	var next int32
	if size == 0 {
		next = int32(len(data))
	} else {
		next, _ = l.Get(1)
	}

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
	return (headerSize + l.Count()) * itemSize
}

func (l *List) Count() int32 {
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

func (l *List) Items() iter.Seq2[int32, int32] {
	return func(yield func(int32, int32) bool) {
		for i := range l.size {
			item, ok := l.Item(i)
			if !ok || !yield(i, item) {
				return
			}
		}
	}
}

func (l *List) Pop() (int32, bool) {
	size := l.Count()
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
	size := l.Count()
	if !l.list.Set(size+headerSize, value) {
		return 0, false
	}

	l.setSize(size + 1)
	return size, true
}
