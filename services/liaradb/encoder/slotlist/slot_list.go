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
	next  int16
	list  int16list.Int16List
}

func New(data []byte) SlotList {
	l := int16list.New(data)
	count, _ := l.Get(0)

	sl := SlotList{
		count: count,
		list:  l,
	}

	sl.initNext()

	return sl
}

func (sl *SlotList) Next() int16 { return sl.next }

// TODO: Remove this
func (sl *SlotList) SetNext(next int16) {
	sl.next = next
}

func (*SlotList) position(i int16) int16 {
	return i*tupleSize + headerSize
}

func (sl *SlotList) initNext() {
	// TODO: Fix this cast
	size := int16(sl.list.Length())
	var next = size

	for slot := range sl.Slots() {
		if !slot.IsFilled() {
			return
		}

		if offset := slot.Offset(); offset < next {
			next = offset
		}
	}

	sl.next = next
}

func (sl *SlotList) Reset() {
	count, _ := sl.list.Get(0)
	sl.count = count
}

func (sl *SlotList) Clear() {
	_ = sl.list.Set(0, 0)
	sl.count = 0
	sl.initNext()
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

func (sl *SlotList) setCount(count int16) {
	if sl.list.Set(0, count) {
		sl.count = count
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

func (sl *SlotList) Insert(offset int16, size int16, i int16) (int16, bool) {
	start := sl.position(i)
	end := sl.position(sl.count)

	if ok := sl.list.ShiftRange(start, end, slotSize); !ok {
		return 0, false
	}

	if ok := sl.list.Set(start, offset); !ok {
		return 0, false
	}

	if ok := sl.list.Set(start+1, size); !ok {
		return 0, false
	}

	count := sl.count
	sl.setCount(count + 1)
	return count, true
}

func (sl *SlotList) Pop() (Slot, bool) {
	if sl.count < 1 {
		return Slot{}, false
	}

	slot, ok := sl.Slot(sl.count - 1)
	if !ok {
		return Slot{}, false
	}

	sl.setCount(sl.count - 1)
	return slot, true
}

func (sl *SlotList) Push(offset int16, size int16) (int16, bool) {
	pos := sl.position(sl.count)
	if !sl.list.Set(pos, offset) {
		return 0, false
	}

	if !sl.list.Set(pos+1, size) {
		return 0, false
	}

	count := sl.count
	sl.setCount(count + 1)
	return count, true
}
