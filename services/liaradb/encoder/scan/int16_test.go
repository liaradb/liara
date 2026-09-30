package scan

import "testing"

func TestInt16(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want int16 = 12345
	_ = SetInt16(data, want)

	if v, _ := Int16(data); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestUint16(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want uint16 = 12345
	_ = SetUint16(data, want)

	if v, _ := Uint16(data); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestIn16__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	b0 := SetInt16(data, 0)

	if l := len(b0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	_, b1 := Int16(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}

func TestUin16__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	b0 := SetUint16(data, 0)

	if l := len(b0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	_, b1 := Uint16(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
