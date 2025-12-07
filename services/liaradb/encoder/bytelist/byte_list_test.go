package bytelist

import "testing"

func TestByteList(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 16))

	b, ok := l.Slice(0, 16)
	if !ok {
		t.Error("should get a buffer")
	}

	if l := len(b); l != 16 {
		t.Errorf("incorrect length: %v, expected: %v", l, 16)
	}

	if _, ok := l.Slice(16, 16); ok {
		t.Error("should not get a buffer starting beyond length")
	}

	if _, ok := l.Slice(0, 20); ok {
		t.Error("should not get a buffer ending beyond length")
	}
}
