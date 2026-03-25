package scan

import "testing"

func TestInt8(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want int8 = 123
	_, ok := SetInt8(data, want)
	if !ok {
		t.Error("unable to write")
	}

	if v, _, ok := Int8(data); !ok {
		t.Error("unable to read")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestByte(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want byte = 123
	_, ok := SetByte(data, want)
	if !ok {
		t.Error("unable to write")
	}

	if v, _, ok := Byte(data); !ok {
		t.Error("unable to read")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestInt8__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 2)
	b0, ok := SetInt8(data, 0)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(b0); l != 1 {
		t.Errorf("incorrect length: %v, expected: %v", l, 1)
	}

	_, b1, ok := Int8(b0)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}

func TestByte__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 2)
	b0, ok := SetByte(data, 0)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(b0); l != 1 {
		t.Errorf("incorrect length: %v, expected: %v", l, 1)
	}

	_, b1, ok := Byte(b0)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
