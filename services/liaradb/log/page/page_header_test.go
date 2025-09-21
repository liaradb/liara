package page

import (
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestPageHeader(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()
	pid := PageID(1)
	tlid := TimeLineID(2)
	rem := RecordLength(3)

	ph := NewPageHeader(pid, tlid, rem)

	if err := ph.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
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
