package record

import (
	"io"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var tm Time = NewTime(time.UnixMicro(1234567890))
	if err := tm.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := tm.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var tm2 Time
	if err := tm2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if tm != tm2 {
		t.Errorf("incorrect value: %v, expected: %v", tm2, tm)
	}
}
