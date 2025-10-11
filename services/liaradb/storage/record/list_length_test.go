package record

import (
	"io"
	"testing"
)

func TestListLength(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var ll ListLength = 123456
	if err := ll.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := ll.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var ll2 ListLength
	if err := ll2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if ll != ll2 {
		t.Errorf("incorrect value: %v, expected: %v", ll2, ll)
	}
}
