package record

import (
	"testing"

	"github.com/liaradb/liaradb/util/testing/iotesting"
)

func TestLogData(t *testing.T) {
	t.Parallel()

	r, w := iotesting.NewReaderWriter()

	want := "abcde"
	ld := NewLogData([]byte(want))

	if l := ld.Length(); l != len(want) {
		t.Errorf("incorrect length: %v, expected: %v", l, len(want))
	}

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

func TestLogData_Compare(t *testing.T) {
	t.Parallel()

	ld0 := NewLogData([]byte("abcde"))
	ld1 := NewLogData([]byte("abcde"))
	ld2 := NewLogData([]byte("12345"))

	if !ld0.Compare(ld1) {
		t.Error("should be equal")
	}

	if ld1.Compare(ld2) {
		t.Error("should not be equal")
	}
}
