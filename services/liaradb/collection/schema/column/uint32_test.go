package column

import "testing"

func TestUInt32(t *testing.T) {
	name := "name"
	u := NewUInt32(name)

	if n := u.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := u.Size(); s != 4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4)
	}

	if tp := u.Type(); tp != TypeUInt32 {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeUInt32)
	}
}
