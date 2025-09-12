package log

import (
	"io"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestLogPageIDTest(t *testing.T) {
	r, w := assert.NewReaderWriter()

	var lpid LogPageID = 123456
	if err := lpid.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var lpid2 LogPageID
	if err := lpid2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if lpid != lpid2 {
		t.Errorf("incorrect value: %v, expected: %v", lpid2, lpid)
	}
}
