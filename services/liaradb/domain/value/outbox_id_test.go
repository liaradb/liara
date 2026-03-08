package value

import (
	"io"
	"testing"

	"github.com/google/uuid"
)

func TestOutboxID(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var e OutboxID = NewOutboxID()
	if err := e.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := e.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var e2 OutboxID
	if err := e2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	s1, s2 := e.String(), e2.String()
	if s1 != s2 {
		t.Errorf("incorrect value: %v, expected: %v", s2, s1)
	}
}

func TestOutboxID_NewOutboxIDFromString(t *testing.T) {
	t.Parallel()

	value := uuid.NewString()

	e, err := NewOutboxIDFromString(value)
	if err != nil {
		t.Error(err)
	}

	if s := e.String(); s != value {
		t.Errorf("incorrect string: %v, expected: %v", s, value)
	}

	_, err = NewOutboxIDFromString("abcde")
	if err == nil {
		t.Error("should return error")
	}
}
