package schema

import "testing"

func TestUInt16Column(t *testing.T) {
	name := "name"
	bc := NewUInt16Column(name)

	if n := bc.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := bc.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if tp := bc.Type(); tp != ColumnTypeUInt16 {
		t.Errorf("incorrect type: %v, expected: %v", tp, ColumnTypeUInt16)
	}
}
