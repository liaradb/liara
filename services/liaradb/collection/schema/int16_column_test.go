package schema

import "testing"

func TestInt16Column(t *testing.T) {
	name := "name"
	bc := NewInt16Column(name)

	if n := bc.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := bc.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if tp := bc.Type(); tp != ColumnTypeInt16 {
		t.Errorf("incorrect type: %v, expected: %v", tp, ColumnTypeInt16)
	}
}
