package link

import "testing"

func TestFileName(t *testing.T) {
	n := "testfile"
	fn := NewFileName(n)

	if s := fn.String(); s != n {
		t.Errorf("incorrect string: %v, expected: %v", s, n)
	}

	want := NewFileBlockID(fn, 1)
	if fbid := fn.FileBlockID(1); fbid != want {
		t.Errorf("incorrect block id: %v, expected: %v", fbid, want)
	}
}
