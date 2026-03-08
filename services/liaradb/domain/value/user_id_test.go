package value

import (
	"io"
	"testing"
)

func TestUserID(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var a = NewUserID("name")
	if err := a.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := a.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var a2 UserID
	if err := a2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if a2 != a {
		t.Errorf("incorrect value: %v, expected: %v", a2, a)
	}
}
