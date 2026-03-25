package scan

import "testing"

func TestInt16(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want int16 = 12345
	_, ok := SetInt16(data, want)
	if !ok {
		t.Error("unable to write")
	}

	if v, _, ok := Int16(data); !ok {
		t.Error("unable to read")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestUint16(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want uint16 = 12345
	_, ok := SetUint16(data, want)
	if !ok {
		t.Error("unable to write")
	}

	if v, _, ok := Uint16(data); !ok {
		t.Error("unable to read")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestIn16__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	b0, ok := SetInt16(data, 0)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(b0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	_, b1, ok := Int16(b0)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}

func TestUin16__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	b0, ok := SetUint16(data, 0)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(b0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	_, b1, ok := Uint16(b0)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
