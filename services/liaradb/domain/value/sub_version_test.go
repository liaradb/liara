package value

import (
	"io"
	"testing"
)

func TestSubVersion(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var gv SubVersion = NewSubVersion(1)
	if err := gv.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := gv.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var gv2 SubVersion
	if err := gv2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	s1, s2 := gv.String(), gv2.String()
	if s1 != s2 {
		t.Errorf("incorrect value: %v, expected: %v", s2, s1)
	}
}

func TestSubVersion_Increment(t *testing.T) {
	t.Parallel()

	sv0 := NewSubVersion(1)
	sv1 := sv0.Increment()
	sv2 := sv1.Increment()

	if v := sv0.Value(); v != 1 {
		t.Errorf("incorrect value: %v, expected: %v", v, 1)
	}
	if v := sv1.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}
	if v := sv2.Value(); v != 3 {
		t.Errorf("incorrect value: %v, expected: %v", v, 3)
	}
}
