package value

import (
	"io"
	"testing"
)

func TestPartitionID(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var e PartitionID = NewPartitionID(1)
	if err := e.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := e.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var e2 PartitionID
	if err := e2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	s1, s2 := e.String(), e2.String()
	if s1 != s2 {
		t.Errorf("incorrect value: %v, expected: %v", s2, s1)
	}
}
