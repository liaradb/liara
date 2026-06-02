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

	sl := SlotList{
		count: count,
		list:  l,
	}

	return sl
}

func (*SlotList) position(i int16) int16 {
	return i*tupleSize + headerSize
}

func (sl *SlotList) Last() (Slot, bool) {
	return sl.Slot(sl.count - 1)
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

func (sl *SlotList) NextSize() int16 {
	return sl.position(sl.count+1) * slotSize
}

func (sl *SlotList) Count() int16 {
	return sl.count
}

func (sl *SlotList) setCount(count int16) {
	if sl.list.Set(0, count) {
		sl.count = count
	}
}

func (sl *SlotList) Slot(i int16) (Slot, bool) {
	if i < 0 || i >= sl.count {
		return Slot{}, false
	}

	pos := sl.position(i)

	a, b, ok := sl.getSlot(pos)
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

func (sl *SlotList) Insert(offset int16, size int16, i int16) (int16, bool) {
	start := sl.position(i)
	end := sl.position(sl.count)

	if ok := sl.list.ShiftRange(start, end, slotSize); !ok {
		return 0, false
	}

	if !sl.setSlot(start, offset, size) {
		return 0, false
	}

	count := sl.count
	sl.setCount(count + 1)
	return count, true
}

func (sl *SlotList) Pop() (Slot, bool) {
	slot, ok := sl.Slot(sl.count - 1)
	if !ok {
		return Slot{}, false
	}

	sl.setCount(sl.count - 1)
	return slot, true
}

func (sl *SlotList) Push(offset int16, size int16) (int16, bool) {
	pos := sl.position(sl.count)
	if !sl.setSlot(pos, offset, size) {
		return 0, false
	}

	count := sl.count
	sl.setCount(count + 1)
	return count, true
}

func (sl *SlotList) getSlot(pos int16) (int16, int16, bool) {
	offset, ok := sl.list.Get(pos)
	if !ok {
		return 0, 0, false
	}

	size, ok := sl.list.Get(pos + 1)
	return offset, size, ok
}

func (sl *SlotList) setSlot(pos, offset, size int16) bool {
	return sl.list.Set(pos, offset) &&
		sl.list.Set(pos+1, size)
}
