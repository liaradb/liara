package slotlist

import (
	"slices"
	"testing"
)

func TestSlot_Range(t *testing.T) {
	t.Parallel()

	data := func(size int) []byte {
		data := make([]byte, size)
		for i := range size {
			data[i] = byte(i)
		}
		return data
	}(110)

	// 	TODO: Should this use a constructor as the fields are private?
	for message, c := range map[string]struct {
		skip  bool
		s     Slot
		start int16
		end   int16
	}{
		"should handle empty slot": {},
		"should handle slot with offset but no size": {
			s:     NewSlot(10, 0),
			start: 10,
			end:   10},
		"should handle slot with size, but no offset": {
			s:     NewSlot(0, 100),
			start: 0,
			end:   100},
		"should handle slot with offset and size": {
			s:     NewSlot(10, 100),
			start: 10,
			end:   110},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			want := data[c.start:c.end]
			if slot := c.s.SliceUnsafe(data); !slices.Equal(slot, want) {
				t.Errorf("incorrect slot: %v, expected: %v", slot, want)
			}
		})
	}
}

func TestSlot_Slice(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)

	s := NewSlot(0, 16)
	b, ok := s.Slice(data)
	if !ok {
		t.Error("should get a buffer")
	}

	if l := len(b); l != 16 {
		t.Errorf("incorrect length: %v, expected: %v", l, 16)
	}

	s = NewSlot(16, 16)
	if _, ok := s.Slice(data); ok {
		t.Error("should not get a buffer starting beyond length")
	}

	s = NewSlot(0, 20)
	if _, ok := s.Slice(data); ok {
		t.Error("should not get a buffer ending beyond length")
	}
}

func TestSlot_Slice__Empty(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)

	s := NewSlot(2, 0)
	b, ok := s.Slice(data)
	if !ok {
		t.Error("should get a buffer")
	}

	if l := len(b); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 16)
	}
}
