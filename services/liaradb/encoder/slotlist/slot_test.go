package slotlist

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/util/testing/should"
)

// TODO: How do we test without using private constructor?
func TestSlot_SliceUnsafe(t *testing.T) {
	t.Parallel()

	data := func(size int) []byte {
		data := make([]byte, size)
		for i := range size {
			data[i] = byte(i)
		}
		return data
	}(110)

	for message, c := range map[string]struct {
		skip  bool
		s     Slot
		start int16
		end   int16
	}{
		"should handle empty slot": {},
		"should handle slot with offset but no size": {
			s:     newSlot(0, 10, 0, data),
			start: 10,
			end:   10},
		"should handle slot with size, but no offset": {
			s:     newSlot(0, 0, 100, data),
			start: 0,
			end:   100},
		"should handle slot with offset and size": {
			s:     newSlot(0, 10, 100, data),
			start: 10,
			end:   110},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			want := data[c.start:c.end]
			if slot := c.s.SliceUnsafe(); !slices.Equal(slot, want) {
				t.Errorf("incorrect slot: %v, expected: %v", slot, want)
			}
		})
	}
}

func TestSlot_Slice(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)

	s := newSlot(0, 0, 16, data)
	b := s.Slice()

	if l := len(b); l != 16 {
		t.Errorf("incorrect length: %v, expected: %v", l, 16)
	}

	s = newSlot(0, 16, 16, data)
	should.Panic(t, "should not get a buffer starting beyond length", func() {
		_ = s.Slice()
	})

	s = newSlot(0, 0, 20, data)
	should.Panic(t, "should not get a buffer ending beyond length", func() {
		_ = s.Slice()
	})
}

func TestSlot_Slice__Empty(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)

	s := newSlot(0, 2, 0, data)
	b := s.Slice()

	if l := len(b); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 16)
	}
}
