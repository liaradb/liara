package column

import "testing"

func TestInt64(t *testing.T) {
	name := "name"
	i := NewInt64(name)

	if n := i.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := i.Size(); s != 8 {
		t.Errorf("incorrect size: %v, expected: %v", s, 8)
	}

	if tp := i.Type(); tp != TypeInt64 {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeInt64)
	}
}
