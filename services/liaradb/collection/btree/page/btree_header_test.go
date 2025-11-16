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

func TestBTreePage_ParentID(t *testing.T) {
	t.Parallel()

	p, _ := newHeader(make([]byte, headerSize))
	if id := p.ParentID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.setParentID(1)
	if id := p.ParentID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}

func TestBTreePage_PrevID(t *testing.T) {
	t.Parallel()

	p, _ := newHeader(make([]byte, headerSize))
	if id := p.PrevID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.setPrevID(1)
	if id := p.PrevID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}

func TestBTreePage_NextID(t *testing.T) {
	t.Parallel()

	p, _ := newHeader(make([]byte, headerSize))
	if id := p.NextID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.setNextID(1)
	if id := p.NextID(); id != 1 {
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
