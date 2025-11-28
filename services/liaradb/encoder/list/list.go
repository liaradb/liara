package list

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/int16list"
)

const (
	headerSize = 2
	itemSize   = 2
)

type List struct {
	size int16
	next int16
	list int16list.Int16List
}

func New(data []byte) List {
	l := int16list.New(data)
	size, _ := l.Get(0)

	var next int16
	if size == 0 {
		next = int16(len(data))
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

func (l *List) Size() int16 {
	return (headerSize + l.Count()) * itemSize
}

func (l *List) Count() int16 {
	return l.size
}

func (l *List) setSize(size int16) {
	if l.list.Set(0, size) {
		l.size = size
	}
}

func (l *List) Next() int16 {
	return l.next
}

func (l *List) SetNext(next int16) {
	if l.list.Set(1, next) {
		l.next = next
	}
}

func (l *List) Item(index int16) (int16, bool) {
	return l.list.Get(index + headerSize)
}

func (l *List) Items() iter.Seq2[int16, int16] {
	return func(yield func(int16, int16) bool) {
		for i := range l.size {
			item, ok := l.Item(i)
			if !ok || !yield(i, item) {
				return
			}
		}
	}
}

func (l *List) Pop() (int16, bool) {
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

func (l *List) Push(value int16) (int16, bool) {
	size := l.Count()
	if !l.list.Set(size+headerSize, value) {
		return 0, false
	}

	l.setSize(size + 1)
	return size, true
}
