package column

import "testing"

func TestUInt16(t *testing.T) {
	name := "name"
	u := NewUInt16(name)

	if n := u.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := u.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if tp := u.Type(); tp != TypeUInt16 {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeUInt16)
	}
}
