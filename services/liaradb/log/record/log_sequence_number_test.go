package record

import (
	"io"
	"testing"
)

func TestLogSequenceNumber(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var lsn LogSequenceNumber = 123456
	if err := lsn.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := lsn.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var lsn2 LogSequenceNumber
	if err := lsn2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if lsn != lsn2 {
		t.Errorf("incorrect value: %v, expected: %v", lsn2, lsn)
	}
}
