package value

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestLogSequenceNumber(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var e EventID = NewEventID()
	if err := e.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := e.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var e2 EventID
	if err := e2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if e != e2 {
		t.Errorf("incorrect value: %v, expected: %v", e2, e)
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
