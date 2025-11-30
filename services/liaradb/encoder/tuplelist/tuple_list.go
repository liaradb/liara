package tuplelist

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/int16list"
)

const (
	headerSize = 1
	itemSize   = 2
	tupleSize  = 2
)

type TupleList struct {
	size int16
	list int16list.Int16List
}

func New(data []byte) TupleList {
	l := int16list.New(data)
	size, _ := l.Get(0)

	return TupleList{
		size: size,
		list: l,
	}
}

func (l *TupleList) Length() int {
	return l.list.Length()
}

func (l *TupleList) Size() int16 {
	return (headerSize + (l.Count() * tupleSize)) * itemSize
}

func (l *TupleList) Count() int16 {
	return l.size
}

func (l *TupleList) setSize(size int16) {
	if l.list.Set(0, size) {
		l.size = size
	}
}

func (l *TupleList) Item(index int16) (int16, int16, bool) {
	a, ok := l.list.Get((tupleSize * index) + headerSize)
	if !ok {
		return 0, 0, false
	}

	b, ok := l.list.Get((tupleSize * index) + 1 + headerSize)
	if !ok {
		return 0, 0, false
	}

	return a, b, true
}

func (l *TupleList) Items() iter.Seq2[int16, int16] {
	return func(yield func(int16, int16) bool) {
		for i := range l.size {
			a, b, ok := l.Item(i)
			if !ok || !yield(a, b) {
				return
			}
		}
	}
}

func (l *TupleList) Insert(a int16, b int16, i int16) (int16, bool) {
	index := i*tupleSize + headerSize

	if ok := l.list.Shift(index, 2); !ok {
		return 0, false
	}

	if ok := l.list.Set(index, a); !ok {
		return 0, false
	}

	if ok := l.list.Set(index+1, b); !ok {
		return 0, false
	}

	size := l.Count()
	l.setSize(size + 1)
	return size, true
}

func (l *TupleList) Pop() (int16, int16, bool) {
	size := l.Count()
	if size < 1 {
		return 0, 0, false
	}

	a, b, ok := l.Item(size - 1)
	if !ok {
		return 0, 0, false
	}

	l.setSize(size - 1)
	return a, b, true
}

func (l *TupleList) Push(a int16, b int16) (int16, bool) {
	size := l.Count()
	if !l.list.Set((tupleSize*size)+headerSize, a) {
		return 0, false
	}

	if !l.list.Set((tupleSize*size)+1+headerSize, b) {
		return 0, false
	}

	l.setSize(size + 1)
	return size, true
}
