package node

import "testing"

func TestBTreeHeader_Level(t *testing.T) {
	t.Parallel()

	h, _ := newHeader(make([]byte, testHeaderSize))
	if l := h.Level(); l != 0 {
		t.Errorf("incorrect value: %v, expected: %v", h, 0)
	}

	h.setLevel(1)
	if l := h.Level(); l != 1 {
		t.Errorf("incorrect value: %v, expected: %v", h, 1)
	}
}

func TestBTreePage_HighID(t *testing.T) {
	t.Parallel()

	h, _ := newHeader(make([]byte, testHeaderSize))
	if id := h.HighID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", h, 0)
	}

	h.SetHighID(1)
	if id := h.HighID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", h, 1)
	}
}

func TestBTreePage_LowID(t *testing.T) {
	t.Parallel()

	h, _ := newHeader(make([]byte, testHeaderSize))
	if id := h.LowID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", h, 0)
	}

	h.SetLowID(1)
	if id := h.LowID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", h, 1)
	}
}
