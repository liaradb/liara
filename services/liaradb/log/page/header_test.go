package page

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/cardboardrobots/assert"
	"github.com/liaradb/liaradb/log/record"
)

func TestHeader(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()
	pid := PageID(1)
	tlid := TimeLineID(2)
	rem := record.Length(3)

	h := NewHeader(pid, tlid, rem)

	if err := h.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := h.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	h2 := Header{}
	if err := h2.Read(r); err != nil {
		t.Fatal(err)
	}

	testHeader(t, h2, pid, tlid, rem)
}

func testHeader(
	t *testing.T,
	h Header,
	pid PageID,
	tlid TimeLineID,
	rem record.Length,
) {
	t.Helper()
	assert.Getter(t, h.ID, pid, "ID")
	assert.Getter(t, h.TimeLineID, tlid, "TimeLineID")
	assert.Getter(t, h.LengthRemaining, rem, "LengthRemaining")
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
