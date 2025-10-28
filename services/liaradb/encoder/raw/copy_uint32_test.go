package raw

import "testing"

func TestCopyUint32_ExactLength(t *testing.T) {
	t.Parallel()

	a := make([]byte, 4)
	var value uint32 = 12345

	if err := CopyUint32(a, value, 0); err != nil {
		t.Error(err)
	}

	result, err := GetUint32(a, 0)
	if err != nil {
		t.Error(err)
	}

	if result != value {
		t.Error("incorrect result")
	}
}

func TestCopyUint32_Underflow(t *testing.T) {
	t.Parallel()

	a := make([]byte, 5)
	var value uint32 = 12345

	if err := CopyUint32(a, value, 1); err != nil {
		t.Error(err)
	}

	result, err := GetUint32(a, 1)
	if err != nil {
		t.Error(err)
	}

	if result != value {
		t.Error("incorrect result")
	}
}

func TestCopyUint32_Overflow(t *testing.T) {
	t.Parallel()

	a := make([]byte, 5)
	var value uint32 = 12345

	if err := CopyUint32(a, value, 2); err == nil {
		t.Error("should not copy")
	}

	if _, err := GetUint32(a, 2); err == nil {
		t.Error("should not get")
	}
}

func TestCopyUint32_BadOffset(t *testing.T) {
	t.Parallel()

	a := make([]byte, 5)
	var value uint32 = 12345

	if err := CopyUint32(a, value, 6); err == nil {
		t.Error("should not copy")
	}

	if err := CopyUint32(a, value, -1); err == nil {
		t.Error("should not copy")
	}

	if _, err := GetUint32(a, 6); err == nil {
		t.Error("should not get")
	}

	if _, err := GetUint32(a, -1); err == nil {
		t.Error("should not get")
	}
}
