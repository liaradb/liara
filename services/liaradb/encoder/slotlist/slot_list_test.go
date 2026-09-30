package slotlist

import (
	"encoding/binary"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/should"
)

type tuple struct {
	a int16
	b int16
}

func TestSlotList_Default(t *testing.T) {
	t.Parallel()

	l := New([]byte{})

	if length := l.Length(); length != 0 {
		t.Errorf("incorrect length: %v, expected: %v", length, 0)
	}

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	should.Panic(t, "should not get value", func() {
		_ = l.Last()
	})
}

func TestSlotList_Push(t *testing.T) {
	t.Parallel()

	data := make([]byte, 32)
	l := New(data)

	if length := l.Length(); length != 32 {
		t.Errorf("incorrect length: %v, expected: %v", length, 32)
	}

	should.Panic(t, "should not get value", func() {
		_ = l.Last()
	})

	if s, i := l.Push(2, 10); i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	} else {
		assertSlot(t, s, 2, 10)
	}

	if s := l.Size(); s != 6 {
		t.Errorf("incorrect size: %v, expected: %v", s, 6)
	}

	if c := l.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	if offset := l.Last().Offset(); offset != 2 {
		t.Errorf("incorrect next: %v, expected: %v", offset, 2)
	}

	if s, i := l.Push(5, 20); i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	} else {
		assertSlot(t, s, 5, 20)
	}

	if s := l.Size(); s != 10 {
		t.Errorf("incorrect size: %v, expected: %v", s, 10)
	}

	if c := l.Count(); c != 2 {
		t.Errorf("incorrect count: %v, expected: %v", c, 2)
	}

	if offset := l.Last().Offset(); offset != 5 {
		t.Errorf("incorrect next: %v, expected: %v", offset, 5)
	}

	slot := l.Slot(0)
	if slot.Offset() != 2 {
		t.Errorf("incorrect value: %v, expected: %v", slot.Offset(), 2)
	}
	if slot.Size() != 10 {
		t.Errorf("incorrect value: %v, expected: %v", slot.Size(), 10)
	}

	slot = l.Slot(1)
	if slot.Offset() != 5 {
		t.Errorf("incorrect value: %v, expected: %v", slot.Offset(), 5)
	}
	if slot.size != 20 {
		t.Errorf("incorrect value: %v, expected: %v", slot.Size(), 20)
	}

	should.Panic(t, "should not have a value", func() {
		_ = l.Slot(2)
	})
}

func TestSlotList_Pop(t *testing.T) {
	t.Parallel()

	data := make([]byte, 18)
	l := New(data)

	s, _ := l.Push(1, 10)
	assertSlot(t, s, 1, 10)

	s, _ = l.Push(2, 20)
	assertSlot(t, s, 2, 20)

	slot := l.Pop()
	if slot.Offset() != 2 {
		t.Errorf("incorrect value: %v, expected: %v", slot.Offset(), 2)
	}
	if slot.Size() != 20 {
		t.Errorf("incorrect value: %v, expected: %v", slot.Size(), 20)
	}

	if s := l.Size(); s != 6 {
		t.Errorf("incorrect size: %v, expected: %v", s, 6)
	}

	if c := l.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	slot = l.Pop()
	if slot.Offset() != 1 {
		t.Errorf("incorrect value: %v, expected: %v", slot.Offset(), 1)
	}
	if slot.Size() != 10 {
		t.Errorf("incorrect value: %v, expected: %v", slot.Size(), 10)
	}

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	should.Panic(t, "should not pop beyond empty", func() {
		_ = l.Pop()
	})
}

func TestSlotList_Replace(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	should.Panic(t, "should not replace non-existant slots", func() {
		l.Replace(10, 10, 0)
	})

	should.Panic(t, "should not replace non-existant slots", func() {
		l.Replace(11, 11, 1)
	})

	_, _ = l.Push(5, 0)

	l.Replace(10, 10, 0)

	s := l.Slot(0)
	if o := s.Offset(); o != 10 {
		t.Errorf("incorrect offset: %v, expected: %v", o, 10)
	}
	if s := s.Size(); s != 10 {
		t.Errorf("incorrect size: %v, expected: %v", s, 10)
	}

	_, _ = l.Push(5, 0)

	l.Replace(11, 11, 1)

	s = l.Slot(1)
	if o := s.Offset(); o != 11 {
		t.Errorf("incorrect offset: %v, expected: %v", o, 11)
	}
	if s := s.Size(); s != 11 {
		t.Errorf("incorrect size: %v, expected: %v", s, 11)
	}
}

func TestSloList_Delete(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	_, _ = l.Push(3, 3)
	_, _ = l.Push(5, 5)

	if l.IsDeleted(0) {
		t.Error("should not be deleted")
	}

	if l.IsDeleted(1) {
		t.Error("should not be deleted")
	}

	l.Delete(0)

	if !l.IsDeleted(0) {
		t.Error("should be deleted")
	}

	if l.IsDeleted(1) {
		t.Error("should not be deleted")
	}

	s := l.Slot(0)
	if o := s.Offset(); o != 0 {
		t.Errorf("incorrect offset: %v, expected: %v", o, 0)
	}
	if s := s.Size(); s != 0 {
		t.Errorf("incorrect size: %v, expected: %v", s, 0)
	}

	s = l.Slot(1)
	if o := s.Offset(); o != 5 {
		t.Errorf("incorrect offset: %v, expected: %v", o, 5)
	}
	if s := s.Size(); s != 5 {
		t.Errorf("incorrect size: %v, expected: %v", s, 5)
	}
}

func TestSlotList_FirstOffset(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	_, _ = l.Push(5, 0)

	if o := l.FirstOffset(); o != 5 {
		t.Fatalf("incorrect offset: %v, expected: %v", o, 5)
	}

	_, _ = l.Push(6, 0)

	if o := l.FirstOffset(); o != 5 {
		t.Fatalf("incorrect offset: %v, expected: %v", o, 5)
	}

	_, _ = l.Push(4, 0)

	if o := l.FirstOffset(); o != 4 {
		t.Fatalf("incorrect offset: %v, expected: %v", o, 4)
	}
}

func TestSlotList_Slots(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	data := []tuple{
		{10, 60},
		{20, 70},
		{30, 80},
		{40, 90},
		{50, 100}}

	for _, i := range data {
		_, _ = l.Push(i.a, i.b)
	}

	result := make([]tuple, 0, len(data))
	for i := range l.Slots() {
		result = append(result, tuple{i.Offset(), i.Size()})
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestSlotList_SlotsReverse(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	data := []tuple{
		{10, 60},
		{20, 70},
		{30, 80},
		{40, 90},
		{50, 100}}

	for _, i := range data {
		_, _ = l.Push(i.a, i.b)
	}

	result := make([]tuple, 0, len(data))
	for i := range l.SlotsReverse() {
		result = append(result, tuple{i.Offset(), i.Size()})
	}

	// Partial iteration
	for range l.SlotsReverse() {
		break
	}

	slices.Reverse(data)
	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestSlotList_SlotsRange(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	data := []tuple{
		{10, 60},
		{20, 70},
		{30, 80},
		{40, 90},
		{50, 100}}

	for _, i := range data {
		_, _ = l.Push(i.a, i.b)
	}

	for message, c := range map[string]struct {
		skip  bool
		want  []tuple
		start link.SlotID
		end   link.SlotID
	}{
		"should iterate the range": {
			want: []tuple{
				{20, 70},
				{30, 80},
				{40, 90}},
			start: 1,
			end:   4,
		},
		"should iterate wrapping the end": {
			want:  data,
			start: 0,
			end:   -1,
		},
		"should iterate wrapping the start": {
			want: []tuple{
				{30, 80},
				{40, 90}},
			start: -4,
			end:   -2,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			result := make([]tuple, 0, len(c.want))
			for i := range l.SlotsRange(c.start, c.end) {
				result = append(result, tuple{i.Offset(), i.Size()})
			}

			// Partial iteration
			for range l.SlotsRange(c.start, c.end) {
				break
			}

			if !slices.Equal(result, c.want) {
				t.Errorf("incorrect result: %v, expected: %v", result, c.want)
			}
		})
	}
}

func TestSlotList_SlotSliceSortedByOffset(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 256))

	data := []tuple{
		{30, 80},
		{20, 70},
		{50, 100},
		{40, 90}}

	want := []tuple{
		{20, 70},
		{30, 80},
		{40, 90},
		{50, 100}}

	for _, d := range data {
		l.Push(d.a, d.b)
	}

	slots := l.SlotSliceSortedByOffset()

	for i, d := range want {
		s := slots[i]
		if o := s.Offset(); o != d.a {
			t.Errorf("incorrect offset: %v, expected: %v", o, d.a)
		}
		if s := s.Size(); s != d.b {
			t.Errorf("incorrect size: %v, expected: %v", s, d.b)
		}
	}
}

func TestSlotList_Insert(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip   bool
		data   []tuple
		want   []tuple
		insert tuple
		index  link.SlotID
	}{
		"should insert into beginning": {
			data: []tuple{
				{20, 70},
				{30, 80},
				{40, 90},
				{50, 100}},
			want: []tuple{
				{10, 60},
				{20, 70},
				{30, 80},
				{40, 90},
				{50, 100}},
			insert: tuple{10, 60},
			index:  0,
		},
		"should insert into middle": {
			data: []tuple{
				{10, 60},
				{20, 70},
				{40, 90},
				{50, 100}},
			want: []tuple{
				{10, 60},
				{20, 70},
				{30, 80},
				{40, 90},
				{50, 100}},
			insert: tuple{30, 80},
			index:  2,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			data := make([]byte, 42)
			l := New(data)

			for _, i := range c.data {
				_, _ = l.Push(i.a, i.b)
			}

			s, _ := l.Insert(c.insert.a, c.insert.b, c.index)
			assertSlot(t, s, c.insert.a, c.insert.b)

			wantCount := link.SlotID(len(c.want))
			if count := l.Count(); count != wantCount {
				t.Errorf("incorrect count: %v, expected: %v", count, wantCount)
			}

			result := make([]tuple, 0, len(c.data))
			for i := range l.Slots() {
				result = append(result, tuple{i.Offset(), i.Size()})
			}

			// Partial iteration
			for range l.Slots() {
				break
			}

			if !slices.Equal(result, c.want) {
				t.Errorf("incorrect result: %v, expected: %v", result, c.want)
			}
		})
	}

	t.Run("should not insert beyond size", func(t *testing.T) {
		t.Parallel()

		data := make([]byte, 6)
		l := New(data)
		for i, slot := range []tuple{
			{10, 60}} {
			s, _ := l.Insert(slot.a, slot.b, link.SlotID(i))
			assertSlot(t, s, slot.a, slot.b)
		}

		should.Panic(t, "should not insert beyond size", func() {
			_, _ = l.Insert(20, 70, 0)
		})
	})
}

func TestSlotList_Reset(t *testing.T) {
	t.Parallel()

	data := make([]byte, 32)

	l := New(data)

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	binary.BigEndian.PutUint16(data, 3)
	l.Reset()

	if s := l.Size(); s != 14 {
		t.Errorf("incorrect size: %v, expected: %v", s, 14)
	}

	if c := l.Count(); c != 3 {
		t.Errorf("incorrect count: %v, expected: %v", c, 3)
	}

	c := 0
	for range l.Slots() {
		c++
	}

	if c != 3 {
		t.Errorf("incorrect count: %v, expected: %v", c, 3)
	}
}

func TestSlotList_Clear(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 16))

	if length := l.Length(); length != 16 {
		t.Errorf("incorrect length: %v, expected: %v", length, 16)
	}

	if _, i := l.Push(1, 2); i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := l.Size(); s != 6 {
		t.Errorf("incorrect size: %v, expected: %v", s, 6)
	}

	if c := l.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	l.Clear()

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	c := 0
	for range l.Slots() {
		c++
	}

	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func assertSlot(t *testing.T, s Slot, offset, size int16) {
	if o := s.Offset(); o != offset {
		t.Errorf("incorrect offset: %v, expected: %v", o, offset)
	}

	if s := s.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}
}
