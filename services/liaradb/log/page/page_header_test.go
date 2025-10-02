package page

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestPageHeader(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()
	pid := PageID(1)
	tlid := TimeLineID(2)
	rem := RecordLength(3)

	ph := NewPageHeader(pid, tlid, rem)

	if err := ph.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := ph.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	ph2 := &PageHeader{}
	if err := ph2.Read(r); err != nil {
		t.Fatal(err)
	}

	testPageHeader(t, ph2, pid, tlid, rem)
}

func testPageHeader(
	t *testing.T,
	ph *PageHeader,
	pid PageID,
	tlid TimeLineID,
	rem RecordLength,
) {
	t.Helper()
	assert.Getter(t, ph.ID, pid, "ID")
	assert.Getter(t, ph.TimeLineID, tlid, "TimeLineID")
	assert.Getter(t, ph.LengthRemaining, rem, "LengthRemaining")
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
