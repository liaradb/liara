package rawpage

import "testing"

func TestRawPage(t *testing.T) {
	p := New(make([]byte, 256))
	if i, _, ok := p.Append(16); !ok {
		t.Error("should get a buffer")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if i, _, ok := p.Append(16); !ok {
		t.Error("should get a buffer")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}
}
