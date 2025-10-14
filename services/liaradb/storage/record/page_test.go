package record

import (
	"io"
	"reflect"
	"testing"

	"github.com/liaradb/liaradb/raw"
)

func TestPage(t *testing.T) {
	t.Parallel()

	const size = 256

	b := raw.NewBuffer(size)

	p := NewPage(size)
	p.Add([]byte{1, 2, 3, 4})
	p.Add([]byte{5, 6, 7, 8})

	b.Clear()
	if err := p.Write(b); err != nil {
		t.Fatal(err)
	}

	if s := p.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	if _, err := b.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	p1 := NewPage(256)

	if err := p1.Read(b); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(p, p1) {
		t.Error("pages do not match")
	}
}
