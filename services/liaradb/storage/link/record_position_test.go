package link

import (
	"io"
	"testing"
)

func TestRecordPosition(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var p RecordPosition = 123
	if err := p.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := p.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var p2 RecordPosition
	if err := p2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if p != p2 {
		t.Errorf("incorrect value: %v, expected: %v", p2, p)
	}

	if s := p.String(); s != "123" {
		t.Errorf("incorrect string: %v, expected: %v", s, "123")
	}
}
