package scan

import "testing"

func TestInt64(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)

	var want int64 = 12345
	_ = SetInt64(data, want)

	if v, _ := Int64(data); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestUint64(t *testing.T) {
	t.Parallel()

	data := make([]byte, 8)

	var want uint64 = 12345
	_ = SetUint64(data, want)

	if v, _ := Uint64(data); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}

func TestIn64__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)
	b0 := SetInt64(data, 0)

	if l := len(b0); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	_, b1 := Int64(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}

func TestUin64__Remainder(t *testing.T) {
	t.Parallel()

	data := make([]byte, 16)
	b0 := SetUint64(data, 0)

	if l := len(b0); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	_, b1 := Uint64(b0)

	if l := len(b1); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
