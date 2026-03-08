package value

import (
	"io"
	"testing"
)

func TestEventName(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var en = NewEventName("name")
	if err := en.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := en.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var en2 EventName
	if err := en2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if en2 != en {
		t.Errorf("incorrect value: %v, expected: %v", en2, en)
	}
}
