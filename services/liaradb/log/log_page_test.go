package log

import (
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestLogPage(t *testing.T) {
	lpid := LogPageID(1)
	tlid := TimeLineID(2)

	lp := NewLogPage(256)
	lp.Init(lpid, tlid)

	r, w := assert.NewReaderWriter()

	if err := lp.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lp2 := &LogPage{}
	if err := lp2.Read(r); err != nil {
		t.Fatal(err)
	}

	assert.Getter(t, lp2.ID, lpid, "ID")
	assert.Getter(t, lp2.TimeLineID, tlid, "TimeLineID")
}

func TestLogPage_Append(t *testing.T) {
	r, w := assert.NewReaderWriter()

	lpid := LogPageID(1)
	tlid := TimeLineID(2)
	lp := NewLogPage(256)
	lp.Init(lpid, tlid)

	data := []byte{1, 2, 3, 4, 5, 6}
	crc := NewCRC(data)

	if err := lp.Append(crc, data); err != nil {
		t.Fatal(err)
	}

	if err := lp.Append(crc, data); err != nil {
		t.Fatal(err)
	}

	if err := lp.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lp2 := NewLogPage(256)
	if err := lp2.Read(r); err != nil {
		t.Fatal(err)
	}

	assert.Getter(t, lp2.ID, lpid, "ID")
	assert.Getter(t, lp2.TimeLineID, tlid, "TimeLineID")
	// TODO: This is not using the public API
	assert.EqualsArray(t, lp.data, lp2.data, "data")

	count := 0
	for _, err := range lp2.Records() {
		count++
		if err != nil {
			t.Fatal(err)
		}
	}

	if count != 2 {
		t.Errorf("incorrect count: %v, expected: %v", count, 2)
	}
}
