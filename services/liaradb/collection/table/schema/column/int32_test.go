package column

import "testing"

func TestInt32(t *testing.T) {
	name := "name"
	i := NewInt32(name)

	if n := i.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := i.Size(); s != 4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4)
	}

	if tp := i.Type(); tp != TypeInt32 {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeInt32)
	}
}
