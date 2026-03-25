package scan

import "testing"

func TestInt32(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want int32 = 12345
	_, ok := SetInt32(data, want)
	if !ok {
		t.Error("unable to write")
	}

	if v, _, ok := Int32(data); !ok {
		t.Error("unable to read")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestUint32(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)

	var want uint32 = 12345
	_, ok := SetUint32(data, want)
	if !ok {
		t.Error("unable to write")
	}

	if v, _, ok := Uint32(data); !ok {
		t.Error("unable to read")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestIn32__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)
	b0, ok := SetInt32(data, 0)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(b0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	_, b1, ok := Int32(b0)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}

func TestUin32__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)
	b0, ok := SetUint32(data, 0)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(b0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	_, b1, ok := Uint32(b0)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
