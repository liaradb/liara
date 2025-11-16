package wrap

import "testing"

func TestInt32(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	i, _ := NewInt32(data)

	var want int32 = 12345
	i.Set(want)

	if v := i.Get(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestInt32_Unsigned(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	i, _ := NewInt32(data)

	var want uint32 = 12345
	i.SetUnsigned(want)

	if v := i.GetUnsigned(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestIn32__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)
	_, b0 := NewInt32(data)

	if l := len(b0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	_, b1 := NewInt32(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
