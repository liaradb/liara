package record

import (
	"io"
	"reflect"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/raw"
)

func TestPage_Add(t *testing.T) {
	t.Parallel()

	const size = 256
	p := NewPage(size)

	items := [][]byte{
		{1, 2, 3, 4},
		{5, 6, 7, 8}}
	for _, i := range items {
		p.Add(i)
	}

	count := 0
	for i, err := range p.Items() {
		if err != nil {
			t.Error(err)
		}
		if !slices.Equal(i, items[count]) {
			t.Errorf("item does not match: %v, expected: %v", i, items[count])
		}
		count++
	}

	if count != len(items) {
		t.Errorf("incorrect count: %v, expected: %v", count, len(items))
	}
}

func TestPage_ReadWrite(t *testing.T) {
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
