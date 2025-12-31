package scan

import "testing"

func TestInt8(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want int8 = 123
	_ = SetInt8(data, want)

	if v, _ := Int8(data); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestByte(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want byte = 123
	_ = SetByte(data, want)

	if v, _ := Byte(data); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestInt8__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 2)
	b0 := SetInt8(data, 0)

	if l := len(b0); l != 1 {
		t.Errorf("incorrect length: %v, expected: %v", l, 1)
	}

	_, b1 := Int8(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}

func TestByte__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 2)
	b0 := SetByte(data, 0)

	if l := len(b0); l != 1 {
		t.Errorf("incorrect length: %v, expected: %v", l, 1)
	}

	_, b1 := Byte(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
