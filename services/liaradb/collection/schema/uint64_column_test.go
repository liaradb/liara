package schema

import "testing"

func TestUInt64Column(t *testing.T) {
	name := "name"
	bc := NewUInt64Column(name)

	if n := bc.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := bc.Size(); s != 8 {
		t.Errorf("incorrect size: %v, expected: %v", s, 8)
	}

	if tp := bc.Type(); tp != ColumnTypeUInt64 {
		t.Errorf("incorrect type: %v, expected: %v", tp, ColumnTypeUInt64)
	}
}
