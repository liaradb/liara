package wrap

import "testing"

func TestInt32(t *testing.T) {
	data := make([]byte, 8)
	i, _ := NewInt32(data)

	var want int32 = 12345
	i.Set(want)

	if v := i.Get(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}
