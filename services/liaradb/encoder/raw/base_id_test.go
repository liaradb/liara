package raw

import "testing"

func TestBaseID__Remainder(t *testing.T) {
	t.Parallel()

	b := NewBaseID()

	data := make([]byte, 20)
	data0 := b.WriteData(data)

	if l := len(data0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	b0 := BaseID{}
	data1 := b0.ReadData(data)

	if l := len(data1); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if b0 != b {
		t.Errorf("incorrect value: %v, expected: %v", b0, b)
	}
}
