package crclist

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/int16list"
	"github.com/liaradb/liaradb/encoder/page"
)

const (
	headerSize = 1
	itemSize   = 2
	tupleSize  = 4
)

type CRCList struct {
	count int16
	list  int16list.Int16List
}

func New(data []byte) CRCList {
	l := int16list.New(data)
	count, _ := l.Get(0)

	return CRCList{
		count: count,
		list:  l,
	}
}

// TODO: Test this
func (l *CRCList) Clear() {
	l.count = 0
}

func (l *CRCList) Length() int {
	return l.list.Length()
}

func (l *CRCList) Size() int16 {
	return (headerSize + (l.Count() * tupleSize)) * itemSize
}

func (l *CRCList) Count() int16 {
	return l.count
}

func (l *CRCList) setSize(size int16) {
	if l.list.Set(0, size) {
		l.count = size
	}
}

func (l *CRCList) Item(index int16) (CRCItem, bool) {
	a, ok := l.list.Get((tupleSize * index) + headerSize)
	if !ok {
		return CRCItem{}, false
	}

	b, ok := l.list.Get((tupleSize * index) + 1 + headerSize)
	if !ok {
		return CRCItem{}, false
	}

	c, ok := l.list.GetInt32((tupleSize * index) + 2 + headerSize)
	if !ok {
		return CRCItem{}, false
	}

	return CRCItem{a, b, page.RestoreCRC(c)}, true
}

type CRCItem struct {
	Offset int16
	Size   int16
	CRC    page.CRC
}

func (l *CRCList) Items() iter.Seq[CRCItem] {
	return func(yield func(CRCItem) bool) {
		for i := range l.count {
			item, ok := l.Item(i)
			if !ok || !yield(item) {
				return
			}
		}
	}
}

func (l *CRCList) ItemsRange(start, end int16) iter.Seq[CRCItem] {
	return func(yield func(CRCItem) bool) {
		if start < 0 {
			start = l.count + 1 + start
		}
		if end < 0 {
			end = l.count + 1 + end
		}
		for i := start; i < end; i++ {
			item, ok := l.Item(i)
			if !ok || !yield(item) {
				return
			}
		}
	}
}

func (l *CRCList) Insert(a int16, b int16, c page.CRC, i int16) (int16, bool) {
	index := i*tupleSize + headerSize
	s := l.count*tupleSize + headerSize

	if ok := l.list.ShiftRange(index, s, 4); !ok {
		return 0, false
	}

	if ok := l.list.Set(index, a); !ok {
		return 0, false
	}

	if ok := l.list.Set(index+1, b); !ok {
		return 0, false
	}

	if ok := l.list.SetInt32(index+2, int32(c.Value())); !ok {
		return 0, false
	}

	size := l.Count()
	l.setSize(size + 1)
	return size, true
}

func (l *CRCList) SetCRC(c page.CRC, i int16) bool {
	index := i*tupleSize + headerSize
	return l.list.SetInt32(index+2, int32(c.Value()))
}

func (l *CRCList) Pop() (CRCItem, bool) {
	size := l.Count()
	if size < 1 {
		return CRCItem{}, false
	}

	item, ok := l.Item(size - 1)
	if !ok {
		return CRCItem{}, false
	}

	l.setSize(size - 1)
	return item, true
}

func (l *CRCList) Push(a int16, b int16, c page.CRC) (int16, bool) {
	size := l.Count()
	if !l.list.Set((tupleSize*size)+headerSize, a) {
		return 0, false
	}

	if !l.list.Set((tupleSize*size)+1+headerSize, b) {
		return 0, false
	}

	if !l.list.SetInt32((tupleSize*size)+2+headerSize, int32(c.Value())) {
		return 0, false
	}

	l.setSize(size + 1)
	return size, true
}
