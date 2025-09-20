package log

import (
	"io"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestRecordLength(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	var rl RecordLength = NewRecordLength([]byte{1, 2, 3, 4, 5})
	if err := rl.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var rl2 RecordLength
	if err := rl2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if rl != rl2 {
		t.Errorf("incorrect value: %v, expected: %v", rl2, rl)
	}
}
