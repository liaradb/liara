package schema

import "testing"

func TestSchema(t *testing.T) {
	id := NewID()
	s := NewSchema(id)

	if i := s.ID(); i != id {
		t.Errorf("incorrect id: %v, expected: %v", i, id)
	}
}

func TestSchema_GetString(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewStringColumn(name, 10))

	_, ok := s.GetString(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetString("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
	}
}
