package column

import "testing"

func TestUInt64(t *testing.T) {
	name := "name"
	u := NewUInt64(name)

	if n := u.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := u.Size(); s != 8 {
		t.Errorf("incorrect size: %v, expected: %v", s, 8)
	}

	if tp := u.Type(); tp != TypeUInt64 {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeUInt64)
	}
}
