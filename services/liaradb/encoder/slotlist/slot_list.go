package slotlist

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/int16list"
)

const (
	headerSize = 1
	itemSize   = 2
	tupleSize  = 2
)

type SlotList struct {
	count int16
	list  int16list.Int16List
}

func New(data []byte) SlotList {
	l := int16list.New(data)
	count, _ := l.Get(0)

	return SlotList{
		count: count,
		list:  l,
	}
}

func (sl *SlotList) Reset() {
	count, _ := sl.list.Get(0)
	sl.count = count
}

func (sl *SlotList) Clear() {
	_ = sl.list.Set(0, 0)
	sl.count = 0
}

func (sl *SlotList) Length() int {
	return sl.list.Length()
}

func (sl *SlotList) Size() int16 {
	return (headerSize + (sl.Count() * tupleSize)) * itemSize
}

func (sl *SlotList) Count() int16 {
	return sl.count
}

func (sl *SlotList) setSize(size int16) {
	if sl.list.Set(0, size) {
		sl.count = size
	}
}

func (sl *SlotList) Item(index int16) (Item, bool) {
	if index >= sl.count {
		return Item{}, false
	}

	a, ok := sl.list.Get((tupleSize * index) + headerSize)
	if !ok {
		return Item{}, false
	}

	b, ok := sl.list.Get((tupleSize * index) + 1 + headerSize)
	if !ok {
		return Item{}, false
	}

	return Item{a, b}, true
}

type Item struct {
	Offset int16
	Size   int16
}

func (sl *SlotList) Items() iter.Seq[Item] {
	return func(yield func(Item) bool) {
		for i := range sl.count {
			item, ok := sl.Item(i)
			if !ok || !yield(item) {
				return
			}
		}
	}
}

func (sl *SlotList) ItemsReverse() iter.Seq[Item] {
	return func(yield func(Item) bool) {
		c := sl.count - 1
		for i := range sl.count {
			item, ok := sl.Item(c - i)
			if !ok || !yield(item) {
				return
			}
		}
	}
}

func (sl *SlotList) ItemsRange(start, end int16) iter.Seq[Item] {
	return func(yield func(Item) bool) {
		if start < 0 {
			start = sl.count + 1 + start
		}
		if end < 0 {
			end = sl.count + 1 + end
		}
		for i := start; i < end; i++ {
			item, ok := sl.Item(i)
			if !ok || !yield(item) {
				return
			}
		}
	}
}

func (sl *SlotList) Insert(a int16, b int16, i int16) (int16, bool) {
	index := i*tupleSize + headerSize
	s := sl.count*tupleSize + headerSize

	if ok := sl.list.ShiftRange(index, s, 2); !ok {
		return 0, false
	}

	if ok := sl.list.Set(index, a); !ok {
		return 0, false
	}

	if ok := sl.list.Set(index+1, b); !ok {
		return 0, false
	}

	size := sl.Count()
	sl.setSize(size + 1)
	return size, true
}

func (sl *SlotList) Pop() (Item, bool) {
	size := sl.Count()
	if size < 1 {
		return Item{}, false
	}

	item, ok := sl.Item(size - 1)
	if !ok {
		return Item{}, false
	}

	sl.setSize(size - 1)
	return item, true
}

func (sl *SlotList) Push(a int16, b int16) (int16, bool) {
	size := sl.Count()
	if !sl.list.Set((tupleSize*size)+headerSize, a) {
		return 0, false
	}

	if !sl.list.Set((tupleSize*size)+1+headerSize, b) {
		return 0, false
	}

	sl.setSize(size + 1)
	return size, true
}
