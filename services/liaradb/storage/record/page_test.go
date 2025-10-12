package record

import (
	"reflect"
	"testing"
)

func TestPage(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	p := NewPage(256)
	p.Add([]byte{1, 2, 3, 4})

	if err := p.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := p.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	p1 := NewPage(256)
	if err := p1.Read(r); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(p, p1) {
		t.Error("pages do not match")
	}
}
