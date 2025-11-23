package wrap

import "testing"

func TestBool(t *testing.T) {
	var b Bool

	b.Set(0)
	if b != 0b00000001 {
		t.Errorf("incorrect value: %v, expected: %v", b, 0b00000001)
	}

	if r := b.Get(0); !r {
		t.Errorf("incorrect result: %v, expected: %v", r, true)
	}
}
