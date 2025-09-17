package log

import (
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestLogPageHeader(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()
	lpid := LogPageID(1)
	tlid := TimeLineID(2)
	rem := LogRecordLength(3)

	lph := newLogPageHeader(lpid, tlid, rem)

	if err := lph.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lph2 := &LogPageHeader{}
	if err := lph2.Read(r); err != nil {
		t.Fatal(err)
	}

	testLogPageHeader(t, lph2, lpid, tlid, rem)
}

func testLogPageHeader(
	t *testing.T,
	lph *LogPageHeader,
	lpid LogPageID,
	tlid TimeLineID,
	rem LogRecordLength,
) {
	t.Helper()
	assert.Getter(t, lph.ID, lpid, "ID")
	assert.Getter(t, lph.TimeLineID, tlid, "TimeLineID")
	assert.Getter(t, lph.LengthRemaining, rem, "LengthRemaining")
}
