package slice

import "testing"

func TestSlice(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)

	b, ok := Slice(data, 0, 16)
	if !ok {
		t.Error("should get a buffer")
	}

	if l := len(b); l != 16 {
		t.Errorf("incorrect length: %v, expected: %v", l, 16)
	}

	if _, ok := Slice(data, 16, 16); ok {
		t.Error("should not get a buffer starting beyond length")
	}

	if _, ok := Slice(data, 0, 20); ok {
		t.Error("should not get a buffer ending beyond length")
	}
}

func TestSlice_Empty(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)

	b, ok := Slice(data, 2, 0)
	if !ok {
		t.Error("should get a buffer")
	}

	if l := len(b); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 16)
	}
}
