package value

import "testing"

func TestVersion_Increment(t *testing.T) {
	version := NewVersion(1)

	if v := version.Value(); v != 1 {
		t.Errorf("incorrect value: %v, expected: %v", v, 1)
	}

	version.Increment()

	if v := version.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}
}
