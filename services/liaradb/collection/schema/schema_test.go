package schema

import "testing"

func TestSchema(t *testing.T) {
	id := NewID()
	s := NewSchema(id)

	if i := s.ID(); i != id {
		t.Errorf("incorrect id: %v, expected: %v", i, id)
	}
}
