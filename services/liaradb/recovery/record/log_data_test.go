package record

import (
	"testing"

	"github.com/liaradb/liaradb/util/testutil"
)

func TestLogData(t *testing.T) {
	t.Parallel()

	r, w := testutil.NewReaderWriter()

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
