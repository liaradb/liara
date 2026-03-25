package scan

import "testing"

func TestInt64(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)

	var want int64 = 12345
	_, ok := SetInt64(data, want)
	if !ok {
		t.Error("unable to set value")
	}

	if v, _, ok := Int64(data); !ok {
		t.Error("unable to get value")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestUint64(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)

	var want uint64 = 12345
	_, ok := SetUint64(data, want)
	if !ok {
		t.Error("unable to set value")
	}

	if v, _, ok := Uint64(data); !ok {
		t.Error("unable to get value")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestIn64__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)
	b0, ok := SetInt64(data, 0)
	if !ok {
		t.Error("unable to set value")
	}

	if l := len(b0); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	_, b1, ok := Int64(b0)
	if !ok {
		t.Error("unable to get value")
	}

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}

func TestUin64__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)
	b0, ok := SetUint64(data, 0)
	if !ok {
		t.Error("unable to set value")
	}

	if l := len(b0); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	_, b1, ok := Uint64(b0)
	if !ok {
		t.Error("unable to get value")
	}

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
