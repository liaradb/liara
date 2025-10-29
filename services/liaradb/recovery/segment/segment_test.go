package segment

import (
	"testing"
)

func TestSegment(t *testing.T) {
	t.Parallel()

	t.Run("should handle default", func(t *testing.T) {
		t.Parallel()

		s := NewSegment(0, 0)

		if v := s.Size(); v != 0 {
			t.Errorf("incorrect size: %v, expected: %v", v, 0)
		}

		if v := s.PageSize(); v != 0 {
			t.Errorf("incorrect page size: %v, expected: %v", v, 0)
		}
	})

	t.Run("should handle values", func(t *testing.T) {
		t.Parallel()

		s := NewSegment(1, 2)

		if v := s.Size(); v != 1 {
			t.Errorf("incorrect size: %v, expected: %v", v, 1)
		}

		if v := s.PageSize(); v != 2 {
			t.Errorf("incorrect page size: %v, expected: %v", v, 2)
		}
	})
}
