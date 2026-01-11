package schema

import "testing"

func TestSchema(t *testing.T) {
	id := NewID()
	s := NewSchema(id)

	if i := s.ID(); i != id {
		t.Errorf("incorrect id: %v, expected: %v", i, id)
	}
}

func TestSchema_GetByte(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewByteColumn(name))

	_, ok := s.GetByte(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetByte("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
	}
}

func TestSchema_GetInt8(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewInt8Column(name))

	_, ok := s.GetInt8(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetInt8("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
	}
}

func TestSchema_GetInt16(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewInt16Column(name))

	_, ok := s.GetInt16(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetInt16("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
	}
}

func TestSchema_GetInt32(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewInt32Column(name))

	_, ok := s.GetInt32(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetInt32("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
	}
}

func TestSchema_GetInt64(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewInt64Column(name))

	_, ok := s.GetInt64(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetInt64("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
	}
}

func TestSchema_GetUInt16(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewUInt16Column(name))

	_, ok := s.GetUInt16(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetUInt16("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
	}
}

func TestSchema_GetUInt32(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewUInt32Column(name))

	_, ok := s.GetUInt32(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetUInt32("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
	}
}

func TestSchema_GetUInt64(t *testing.T) {
	s := NewSchema(NewID())

	name := "name"
	s.AddColumn(NewUInt64Column(name))

	_, ok := s.GetUInt64(name)
	if !ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, true)
	}

	_, ok = s.GetUInt64("other")
	if ok {
		t.Errorf("incorrect result: %v, expected: %v", ok, false)
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
