package log

import (
	"reflect"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestLogPageReader(t *testing.T) {
	r, w := assert.NewReaderWriter()
	lpid, tlid, lp := createPage()

	if err := lp.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lp2 := &LogPageWriter{}
	if err := lp2.Read(r); err != nil {
		t.Fatal(err)
	}

	assert.Getter(t, lp2.ID, lpid, "ID")
	assert.Getter(t, lp2.TimeLineID, tlid, "TimeLineID")
}

func TestLogPageReader_Append(t *testing.T) {
	r, w := assert.NewReaderWriter()
	lpid, tlid, lp := createPage()

	lr, data, err := createRecord()
	if err != nil {
		t.Fatal(err)
	}

	crc := NewCRC(data)

	if err := lp.append(crc, data); err != nil {
		t.Fatal(err)
	}

	if err := lp.append(crc, data); err != nil {
		t.Fatal(err)
	}

	if err := lp.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lpr := newLogPageReader(256)
	p, err := lpr.Read(r)
	if err != nil {
		t.Fatal(err)
	}

	testLogPageHeader(t, p, lpid, tlid, 0)

	count := 0
	for r, err := range lpr.Records() {
		count++
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(r, lr) {
			t.Error("data does not match")
		}
	}

	if count != 2 {
		t.Errorf("incorrect count: %v, expected: %v", count, 2)
	}
}
