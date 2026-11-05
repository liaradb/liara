package column

import "testing"

func TestString(t *testing.T) {
	name := "name"
	size := 10
	s := NewString(name, size)

	if n := s.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := s.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	if tp := s.Type(); tp != TypeString {
		t.Errorf("incorrect type: %v, expected: %v", tp, TypeString)
	}
}
