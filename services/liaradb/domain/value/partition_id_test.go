package value

import (
	"io"
	"testing"
)

func TestPartitionID(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var pid PartitionID = NewPartitionID(1)
	if err := pid.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := pid.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var e2 PartitionID
	if err := e2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	s1, s2 := pid.String(), e2.String()
	if s1 != s2 {
		t.Errorf("incorrect value: %v, expected: %v", s2, s1)
	}
}

func TestPartitionID_Value(t *testing.T) {
	t.Parallel()

	var pid0 PartitionID = NewPartitionID(1)

	if v := pid0.Value(); v != 1 {
		t.Errorf("incorrect value: %v, expected: %v", v, 1)
	}

	var pid1 PartitionID = NewPartitionID(2)

	if v := pid1.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}
}
