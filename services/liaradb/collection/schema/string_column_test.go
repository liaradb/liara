package schema

import "testing"

func TestStringColumn(t *testing.T) {
	name := "name"
	size := 10
	bc := NewStringColumn(name, size)

	if n := bc.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := bc.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	if tp := bc.Type(); tp != ColumnTypeString {
		t.Errorf("incorrect type: %v, expected: %v", tp, ColumnTypeString)
	}
}
