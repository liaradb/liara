package page

import "testing"

func TestPage(t *testing.T) {
	t.Parallel()

	p := NewPage(32, 4, 4)
	header, data := p.Next(8)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	i, ok := p.Commit(8)
	if !ok {
		t.Error("should commit")
	}

	header, data, ok = p.Slot(i)
	if !ok {
		t.Error("should get slot")
	}

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}
}
