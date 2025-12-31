package record

import (
	"io"
	"testing"

	"github.com/liaradb/liaradb/util/testutil"
)

func TestLength(t *testing.T) {
	t.Parallel()

	r, w := testutil.NewReaderWriter()

	var l Length = NewLength(5)
	if err := l.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var l2 Length
	if err := l2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if l != l2 {
		t.Errorf("incorrect value: %v, expected: %v", l2, l)
	}
}
