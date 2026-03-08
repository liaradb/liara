package value

import (
	"io"
	"testing"
)

func TestRowName(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var rn = NewRowName("name")
	if err := rn.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := rn.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var rn2 RowName
	if err := rn2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if rn2 != rn {
		t.Errorf("incorrect value: %v, expected: %v", rn2, rn)
	}
}
