package wrap

import "testing"

func TestInt16(t *testing.T) {
	t.Parallel()

	data := make([]byte, 2)
	i, _ := NewInt16(data)

	var want int16 = 12345
	i.Set(want)

	if v := i.Get(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestInt16_Unsigned(t *testing.T) {
	t.Parallel()

	data := make([]byte, 2)
	i, _ := NewInt16(data)

	var want uint16 = 12345
	i.SetUnsigned(want)

	if v := i.GetUnsigned(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestInt16__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	_, b0 := NewInt16(data)

	if l := len(b0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	_, b1 := NewInt16(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
