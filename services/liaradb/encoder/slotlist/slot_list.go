package slotlist

import (
	"iter"

	"github.com/liaradb/liaradb/encoder/int16list"
)

const (
	headerSize = 1
	slotSize   = 2
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

func (*SlotList) position(i int16) int16 {
	return i*tupleSize + headerSize
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
	return sl.position(sl.count) * slotSize
}

func (sl *SlotList) Count() int16 {
	return sl.count
}

func (sl *SlotList) setSize(size int16) {
	if sl.list.Set(0, size) {
		sl.count = size
	}
}

func (sl *SlotList) Slot(i int16) (Slot, bool) {
	if i >= sl.count {
		return Slot{}, false
	}

	pos := sl.position(i)

	a, ok := sl.list.Get(pos)
	if !ok {
		return Slot{}, false
	}

	b, ok := sl.list.Get(pos + 1)
	if !ok {
		return Slot{}, false
	}

	return Slot{a, b}, true
}

func (sl *SlotList) Slots() iter.Seq[Slot] {
	return func(yield func(Slot) bool) {
		for i := range sl.count {
			slot, ok := sl.Slot(i)
			if !ok || !yield(slot) {
				return
			}
		}
	}
}

func (sl *SlotList) SlotsReverse() iter.Seq[Slot] {
	return func(yield func(Slot) bool) {
		c := sl.count - 1
		for i := range sl.count {
			slot, ok := sl.Slot(c - i)
			if !ok || !yield(slot) {
				return
			}
		}
	}
}

func (sl *SlotList) SlotsRange(start, end int16) iter.Seq[Slot] {
	return func(yield func(Slot) bool) {
		if start < 0 {
			start = sl.count + 1 + start
		}
		if end < 0 {
			end = sl.count + 1 + end
		}
		for i := start; i < end; i++ {
			slot, ok := sl.Slot(i)
			if !ok || !yield(slot) {
				return
			}
		}
	}
}

func (sl *SlotList) Insert(a int16, b int16, i int16) (int16, bool) {
	start := sl.position(i)
	end := sl.position(sl.count)

	if ok := sl.list.ShiftRange(start, end, slotSize); !ok {
		return 0, false
	}

	if ok := sl.list.Set(start, a); !ok {
		return 0, false
	}

	if ok := sl.list.Set(start+1, b); !ok {
		return 0, false
	}

	size := sl.Count()
	sl.setSize(size + 1)
	return size, true
}

func (sl *SlotList) Pop() (Slot, bool) {
	size := sl.Count()
	if size < 1 {
		return Slot{}, false
	}

	slot, ok := sl.Slot(size - 1)
	if !ok {
		return Slot{}, false
	}

	sl.setSize(size - 1)
	return slot, true
}

func (sl *SlotList) Push(a int16, b int16) (int16, bool) {
	size := sl.Count()
	pos := sl.position(size)
	if !sl.list.Set(pos, a) {
		return 0, false
	}

	if !sl.list.Set(pos+1, b) {
		return 0, false
	}

	sl.setSize(size + 1)
	return size, true
}
