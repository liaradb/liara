package record

import (
	"io"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestLength(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	var l Length = NewLength([]byte{1, 2, 3, 4, 5})
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
