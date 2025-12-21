package wrap

import "testing"

func TestBool_Defaults(t *testing.T) {
	t.Parallel()

	var b Bool
	testBool(t, b, [8]bool{})
}

func TestBool(t *testing.T) {
	t.Parallel()

	var b Bool

	b.Set(0, true)
	testBool(t, b, [8]bool{true})

	b.Set(0, false)
	testBool(t, b, [8]bool{})

	b.Set(3, true)
	b.Set(5, true)
	testBool(t, b, [8]bool{false, false, false, true, false, true})

	b.Set(3, false)
	testBool(t, b, [8]bool{false, false, false, false, false, true})
}

func testBool(t *testing.T, b Bool, values [8]bool) {
	t.Helper()
	for i, v := range values {
		if r := b.Get(byte(i)); r != v {
			t.Errorf("incorrect result: (%v, %v), expected: (%v, %v)", i, r, i, v)
		}
	}
}
