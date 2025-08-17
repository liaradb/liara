package raw

import "testing"

func TestCopyString_ExactLength(t *testing.T) {
	t.Parallel()

	a := make([]byte, 8)
	value := "abc"

	if err := CopyString(a, value, 0); err != nil {
		t.Error(err)
	}

	result, err := GetString(a, 0)
	if err != nil {
		t.Error("should get")
	}

	if result != value {
		t.Error("incorrect result")
	}
}
