package value

import "testing"

func TestLimit(t *testing.T) {
	l := Limit(1)

	if v := l.Value(); v != 1 {
		t.Errorf("incorrect value: %v, expected: %v", v, 1)
	}
}
