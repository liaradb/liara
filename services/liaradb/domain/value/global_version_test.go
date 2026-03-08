package value

import (
	"io"
	"testing"
)

func TestGlobalVersion(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var gv GlobalVersion = NewGlobalVersion(1)
	if err := gv.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := gv.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var gv2 GlobalVersion
	if err := gv2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	s1, s2 := gv.String(), gv2.String()
	if s1 != s2 {
		t.Errorf("incorrect value: %v, expected: %v", s2, s1)
	}
}
