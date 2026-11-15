package column

import "testing"

func TestInt16(t *testing.T) {
	name := "name"
	i := NewInt16(name)

	if n := i.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := i.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if tp := i.Type(); tp != TypeInt16 {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeInt16)
	}
}
