package schema

import "testing"

func TestInt32Column(t *testing.T) {
	name := "name"
	bc := NewInt32Column(name)

	if n := bc.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := bc.Size(); s != 4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4)
	}

	if tp := bc.Type(); tp != ColumnTypeInt32 {
		t.Errorf("incorrect type: %v, expected: %v", tp, ColumnTypeInt32)
	}
}
