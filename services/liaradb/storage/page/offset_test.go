package page

import (
	"io"
	"testing"
)

func TestOffset(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var o Offset = 123456
	if err := o.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := o.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var o2 Offset
	if err := o2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if o != o2 {
		t.Errorf("incorrect value: %v, expected: %v", o2, o)
	}
}
