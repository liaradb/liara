package schema

import "testing"

func TestInt8Column(t *testing.T) {
	name := "name"
	bc := NewInt8Column(name)

	if n := bc.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := bc.Size(); s != 1 {
		t.Errorf("incorrect size: %v, expected: %v", s, 1)
	}

	if tp := bc.Type(); tp != ColumnTypeInt8 {
		t.Errorf("incorrect type: %v, expected: %v", tp, ColumnTypeInt8)
	}
}
