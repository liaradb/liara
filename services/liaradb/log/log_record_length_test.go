package log

import (
	"io"
	"testing"
)

func TestLogRecordLength(t *testing.T) {
	r, w := createReaderWriter()

	var lrl LogRecordLength = NewLogRecordLength([]byte{1, 2, 3, 4, 5})
	if err := lrl.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var lrl2 LogRecordLength
	if err := lrl2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if lrl != lrl2 {
		t.Errorf("incorrect value: %v, expected: %v", lrl2, lrl)
	}
}
