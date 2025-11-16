package wrap

import "testing"

func TestInt64(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)
	i, _ := NewInt64(data)

	var want int64 = 12345
	i.Set(want)

	if v := i.Get(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestInt64_Unsigned(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)
	i, _ := NewInt64(data)

	var want uint64 = 12345
	i.SetUnsigned(want)

	if v := i.GetUnsigned(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestIn64__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)
	_, b0 := NewInt64(data)

	if l := len(b0); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	_, b1 := NewInt64(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
