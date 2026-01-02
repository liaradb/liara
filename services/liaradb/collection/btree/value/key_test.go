package value

import "testing"

func TestKey(t *testing.T) {
	a := NewKey([]byte("a"))
	b := NewKey([]byte("b"))

	if gt := a.Greater(b); gt {
		t.Error("should not be greater")
	}
}
