package link

import "testing"

func TestFileName(t *testing.T) {
	n := "testfile"
	fn := NewFileName(n)

	if s := fn.String(); s != n {
		t.Errorf("incorrect string: %v, expected: %v", s, n)
	}

	want := NewBlockID(n, 1)
	if bid := fn.BlockID(1); bid != want {
		t.Errorf("incorrect block id: %v, expected: %v", bid, want)
	}
}
