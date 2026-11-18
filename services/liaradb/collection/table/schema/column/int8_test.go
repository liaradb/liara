package column

import "testing"

func TestInt8(t *testing.T) {
	name := "name"
	i := NewInt8(name)

	if n := i.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := i.Size(); s != 1 {
		t.Errorf("incorrect size: %v, expected: %v", s, 1)
	}

	if tp := i.Type(); tp != TypeInt8 {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeInt8)
	}

	if v := i.Value(); v != 0 {
		t.Errorf("incorrect value: %v, expected: %v", v, 0)
	}
}
