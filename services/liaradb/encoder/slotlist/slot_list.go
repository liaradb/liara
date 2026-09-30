package slotlist

import (
	"iter"
	"slices"

	"github.com/liaradb/liaradb/encoder/int16list"
	"github.com/liaradb/liaradb/storage/link"
)

// TODO: Use slices to shrink body as slotlist grows
// Store original body size
// When adding to the slotlist, subslice the body

const (
	headerSize = 1
	slotSize   = 2
	tupleSize  = 2
)

type SlotList struct {
	count link.SlotID
	data  []byte
	list  int16list.Int16List
}

func New(data []byte) SlotList {
	l := int16list.New(data)
	var count int16
	if l.Size() > 0 {
		count = l.Get(0)
	}

	return SlotList{
		count: link.SlotID(count),
		data:  data,
		list:  l,
	}
}

func (*SlotList) position(i link.SlotID) int16 {
	// TODO: Fix this cast
	return int16(i*tupleSize + headerSize)
}

func (sl *SlotList) Last() Slot {
	return sl.Slot(sl.count - 1)
}

func (sl *SlotList) FirstOffset() (offset int16) {
	offset = int16(sl.list.Length())

	for i := range sl.count {
		s := sl.Slot(i)
		if i == 0 || s.offset < offset {
			offset = s.offset
		}
	}

	return
}

func (sl *SlotList) Reset() {
	var count int16
	if sl.list.Size() > 0 {
		count = sl.list.Get(0)
	}
	sl.count = link.SlotID(count)
}

func (sl *SlotList) Clear() {
	if sl.list.Size() > 0 {
		sl.list.Set(0, 0)
	}
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

func (sl *SlotList) Count() link.SlotID {
	return sl.count
}

func (sl *SlotList) setCount(count link.SlotID) {
	sl.list.Set(0, count.Value())
	sl.count = count
}

func (sl *SlotList) Slot(i link.SlotID) Slot {
	if i < 0 || i >= sl.count {
		panic("invalid slot")
	}

	pos := sl.position(i)
	a, b := sl.getSlot(pos)
	return newSlot(i, a, b, sl.data)
}

func (sl *SlotList) Slots() iter.Seq[Slot] {
	return func(yield func(Slot) bool) {
		for i := range sl.count {
			slot := sl.Slot(i)
			if !yield(slot) {
				return
			}
		}
	}
}

func (sl *SlotList) SlotsReverse() iter.Seq[Slot] {
	return func(yield func(Slot) bool) {
		c := sl.count - 1
		for i := range sl.count {
			slot := sl.Slot(c - i)
			if !yield(slot) {
				return
			}
		}
	}
}

func (sl *SlotList) SlotsRange(start, end link.SlotID) iter.Seq[Slot] {
	return func(yield func(Slot) bool) {
		if start < 0 {
			start = sl.count + 1 + start
		}
		if end < 0 {
			end = sl.count + 1 + end
		}
		for i := start; i < end; i++ {
			slot := sl.Slot(i)
			if !yield(slot) {
				return
			}
		}
	}
}

func (sl *SlotList) Insert(offset int16, size int16, i link.SlotID) (Slot, link.SlotID) {
	start := sl.position(i)
	end := sl.position(sl.count)

	sl.list.ShiftRange(start, end, slotSize)
	sl.setSlot(start, offset, size)

	count := sl.count
	sl.setCount(count + 1)
	return newSlot(i, offset, size, sl.data), count
}

func (sl *SlotList) Pop() Slot {
	if sl.count == 0 {
		panic("invalid slot")
	}

	slot := sl.Slot(sl.count - 1)
	sl.setCount(sl.count - 1)
	return slot
}

func (sl *SlotList) Push(offset int16, size int16) (Slot, link.SlotID) {
	pos := sl.position(sl.count)
	sl.setSlot(pos, offset, size)

	count := sl.count
	sl.setCount(count + 1)

	return newSlot(count, offset, size, sl.data), count
}

func (sl *SlotList) Replace(offset, size int16, i link.SlotID) {
	if i >= sl.count {
		panic("invalid slot")
	}

	pos := sl.position(i)
	sl.setSlot(pos, offset, size)
}

func (sl *SlotList) Delete(i link.SlotID) {
	sl.Replace(0, 0, i)
}

func (sl *SlotList) IsDeleted(i link.SlotID) bool {
	return sl.Slot(i).isDeleted()
}

func (sl *SlotList) getSlot(pos int16) (int16, int16) {
	offset := sl.list.Get(pos)
	size := sl.list.Get(pos + 1)
	return offset, size
}

func (sl *SlotList) setSlot(pos, offset, size int16) {
	sl.list.Set(pos, offset)
	sl.list.Set(pos+1, size)
}

func (sl *SlotList) SlotSliceSortedByOffset() []Slot {
	slots := sl.slotSlice()
	slices.SortFunc(slots, func(a, b Slot) int { return int(a.offset) - int(b.offset) })
	return slots
}

func (sl *SlotList) slotSlice() []Slot {
	slots := make([]Slot, 0, sl.count)
	for s := range sl.Slots() {
		slots = append(slots, s)
	}
	return slots
}
