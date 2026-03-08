package value

import (
	"io"
	"testing"
)

func TestCorrelationID(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var cid = NewCorrelationID("name")
	if err := cid.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := cid.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var cid2 CorrelationID
	if err := cid2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if cid2 != cid {
		t.Errorf("incorrect value: %v, expected: %v", cid2, cid)
	}
}
