package page

import "testing"

func TestBTreeHeader_Level(t *testing.T) {
	t.Parallel()

	p, _ := newHeader(make([]byte, headerSize))
	if l := p.Level(); l != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.setLevel(1)
	if l := p.Level(); l != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}

func TestBTreePage_HighID(t *testing.T) {
	t.Parallel()

	p, _ := newHeader(make([]byte, headerSize))
	if id := p.HighID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.setHighID(1)
	if id := p.HighID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}

func TestBTreePage_LowID(t *testing.T) {
	t.Parallel()

	p, _ := newHeader(make([]byte, headerSize))
	if id := p.LowID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.setLowID(1)
	if id := p.LowID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}
