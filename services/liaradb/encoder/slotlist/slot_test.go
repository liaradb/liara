package slotlist

import "testing"

func TestSlot_Range(t *testing.T) {
	t.Parallel()

	// 	TODO: Should this use a constructor as the fields are private?
	for message, c := range map[string]struct {
		skip  bool
		s     Slot
		start int16
		end   int16
	}{
		"should handle empty slot": {},
		"should handle slot with offset but no size": {
			s:     Slot{10, 0},
			start: 10,
			end:   10},
		"should handle slot with size, but no offset": {
			s:     Slot{0, 100},
			start: 0,
			end:   100},
		"should handle slot with offset and size": {
			s:     Slot{10, 100},
			start: 10,
			end:   110},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if start, end := c.s.Range(); start != c.start {
				t.Errorf("incorrect start: %v, expected: %v", start, c.start)
			} else if end != c.end {
				t.Errorf("incorrect end: %v, expected: %v", end, c.end)
			}
		})
	}
}
