package schema

import "testing"

func TestByteColumn(t *testing.T) {
	name := "name"
	bc := NewByteColumn(name)

	if n := bc.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}

	if s := bc.Size(); s != 1 {
		t.Errorf("incorrect size: %v, expected: %v", s, 1)
	}

	if tp := bc.Type(); tp != ColumnTypeByte {
		t.Errorf("incorrect type: %v, expected: %v", tp, ColumnTypeByte)
	}
}
