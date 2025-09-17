package log

import (
	"io"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestTransactionID(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	var tid TransactionID = 123456
	if err := tid.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var tid2 TransactionID
	if err := tid2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if tid != tid2 {
		t.Errorf("incorrect value: %v, expected: %v", tid2, tid)
	}
}
