package scan

import "testing"

func TestInt32(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want int32 = 12345
	_ = SetInt32(data, want)

	if v, _ := Int32(data); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestUint32(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want uint32 = 12345
	_ = SetUint32(data, want)

	if v, _ := Uint32(data); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestIn32__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)
	b0 := SetInt32(data, 0)

	if l := len(b0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	_, b1 := Int32(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}

func TestUin32__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)
	b0 := SetUint32(data, 0)

	if l := len(b0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	_, b1 := Uint32(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
