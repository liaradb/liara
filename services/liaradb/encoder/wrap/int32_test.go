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

func TestInt32_Unsigned(t *testing.T) {
	data := make([]byte, 8)
	i, _ := NewInt32(data)

	var want uint32 = 12345
	i.SetUnsigned(want)

	if v := i.GetUnsigned(); v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}
