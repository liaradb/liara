package raw

import "testing"

func TestBaseString(t *testing.T) {
	size := 32
	n := BaseString("name")
	data := make([]byte, size)
	_ = n.WriteData(data, size)

	var r BaseString
	r.ReadData(data, size)
	if r.String() != n.String() {
		t.Errorf("incorrect result: %v, expected: %v", r.String(), n.String())
	}
}
