package column

import "testing"

func TestByte(t *testing.T) {
	name := "name"
	b := NewByte(name)

	if n := b.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := b.Size(); s != 1 {
		t.Errorf("incorrect size: %v, expected: %v", s, 1)
	}

	if tp := b.Type(); tp != TypeByte {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeByte)
	}
}
