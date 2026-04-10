package record

import (
	"io"
	"testing"
)

func TestLogSequenceNumber(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var lsn = NewLogSequenceNumber(123456)
	if err := lsn.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := lsn.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var lsn2 LogSequenceNumber
	if err := lsn2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if lsn != lsn2 {
		t.Errorf("incorrect value: %v, expected: %v", lsn2, lsn)
	}
}

func TestLogSequenceNumber_Increment_Decrement(t *testing.T) {
	t.Parallel()

	lsn0 := NewLogSequenceNumber(0)

	lsn1 := lsn0.Increment()
	want1 := lsn0.Value() + 1
	if v := lsn1.Value(); v != want1 {
		t.Errorf("incorrect value: %v, expected: %v", v, want1)
	}

	lsn2 := lsn1.Increment()
	want2 := lsn1.Value() + 1
	if v := lsn2.Value(); v != want2 {
		t.Errorf("incorrect value: %v, expected: %v", v, want2)
	}

	lsn3 := lsn2.Decrement()
	want3 := lsn2.Value() - 1
	if v := lsn3.Value(); v != want3 {
		t.Errorf("incorrect value: %v, expected: %v", v, want3)
	}

	lsn4 := lsn3.Decrement()
	want4 := lsn3.Value() - 1
	if v := lsn4.Value(); v != want4 {
		t.Errorf("incorrect value: %v, expected: %v", v, want4)
	}
}
