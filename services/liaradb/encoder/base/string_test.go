package base

import (
	"slices"
	"testing"
)

func TestString(t *testing.T) {
	size := 32
	n := String("name")
	data := make([]byte, size)
	_ = n.WriteData(data, size)

	var r String
	data0 := r.ReadData(data, size)

	if l := len(data0); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}

	if r.String() != n.String() {
		t.Errorf("incorrect result: %v, expected: %v", r.String(), n.String())
	}

	if b := r.Bytes(); !slices.Equal(b, []byte("name")) {
		t.Errorf("incorrect bytes: %v, expected: %v", b, []byte("name"))
	}

	if l := r.Length(); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if s := r.Size(); s != 4+4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4+4)
	}
}
