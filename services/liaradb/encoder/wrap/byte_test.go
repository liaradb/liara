package wrap

import "testing"

func TestByte(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	i, _ := NewByte(data)

	var want int8 = 123
	i.Set(want)

	if v := i.Get(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestByte_Unsigned(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	i, _ := NewByte(data)

	var want byte = 123
	i.SetUnsigned(want)

	if v := i.GetUnsigned(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestByte__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 2)
	_, b0 := NewByte(data)

	if l := len(b0); l != 1 {
		t.Errorf("incorrect length: %v, expected: %v", l, 1)
	}

	_, b1 := NewByte(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
