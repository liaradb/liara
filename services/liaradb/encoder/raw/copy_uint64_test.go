package raw

import "testing"

func TestCopyUint64_ExactLength(t *testing.T) {
	t.Parallel()

	a := make([]byte, 8)
	var value uint64 = 12345

	if err := CopyUint64(a, value, 0); err != nil {
		t.Error(err)
	}

	result, err := GetUint64(a, 0)
	if err != nil {
		t.Error(err)
	}

	if result != value {
		t.Error("incorrect result")
	}
}

func TestCopyUint64_Underflow(t *testing.T) {
	t.Parallel()

	a := make([]byte, 9)
	var value uint64 = 12345

	if err := CopyUint64(a, value, 1); err != nil {
		t.Error(err)
	}

	result, err := GetUint64(a, 1)
	if err != nil {
		t.Error(err)
	}

	if result != value {
		t.Error("incorrect result")
	}
}

func TestCopyUint64_Overflow(t *testing.T) {
	t.Parallel()

	a := make([]byte, 5)
	var value uint64 = 12345

	if err := CopyUint64(a, value, 2); err == nil {
		t.Error("should not copy")
	}

	if _, err := GetUint64(a, 2); err == nil {
		t.Error("should not get")
	}
}

func TestCopyUint64_BadOffset(t *testing.T) {
	t.Parallel()

	a := make([]byte, 9)
	var value uint64 = 12345

	if err := CopyUint64(a, value, 6); err == nil {
		t.Error("should not copy")
	}

	if err := CopyUint64(a, value, -1); err == nil {
		t.Error("should not copy")
	}

	if _, err := GetUint64(a, 6); err == nil {
		t.Error("should not get")
	}

	if _, err := GetUint64(a, -1); err == nil {
		t.Error("should not get")
	}
}
