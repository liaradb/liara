package record

import (
	"io"
	"testing"

	"github.com/liaradb/liaradb/util/testing/testutil"
)

func TestTransactionID(t *testing.T) {
	t.Parallel()

	r, w := testutil.NewReaderWriter()

	var tid = NewTransactionID(123456)
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
