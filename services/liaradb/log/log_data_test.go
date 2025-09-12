package log

import (
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestLogData(t *testing.T) {
	r, w := assert.NewReaderWriter()

	want := "abcde"
	ld := NewLogData([]byte(want))

	if err := ld.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	ld2 := LogData{}
	if err := ld2.Read(r); err != nil {
		t.Fatal(err)
	}

	if result := string(ld2.Bytes()); result != want {
		t.Errorf("incorrect value: %v, expected: %v", result, want)
	}
}
