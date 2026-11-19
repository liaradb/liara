package link

import "testing"

func TestFileName(t *testing.T) {
	n := "testfile"
	fn := NewFileName(n)

	if s := fn.String(); s != n {
		t.Errorf("incorrect string: %v, expected: %v", s, n)
	}

	if s, want := fn.Size(), len(n)+4; s != want {
		t.Errorf("incorrect size: %v, expected: %v", s, want)
	}

	want := NewFileBlockID(fn, 1)
	if fbid := fn.FileBlockID(1); fbid != want {
		t.Errorf("incorrect block id: %v, expected: %v", fbid, want)
	}
}

func TestFileName_ReadWrite(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	fn0 := NewFileName("testfile")
	if err := fn0.Write(w); err != nil {
		t.Fatal(err)
	}

	var fn1 FileName
	if err := fn1.Read(r); err != nil {
		t.Fatal(err)
	}

	if fn1 != fn0 {
		t.Errorf("incorrect result: %v, expected: %v", fn1, fn0)
	}
}
