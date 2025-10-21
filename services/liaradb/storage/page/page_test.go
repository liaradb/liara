package page

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
		if err := p.Add(i); err != nil {
			t.Error(err)
		}
	}

	result := make([]Item, 0)

	for i, err := range p.Items() {
		if err != nil {
			t.Error(err)
		}

		result = append(result, i)
	}

	if !slices.EqualFunc(result, items, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, items)
	}
}

func TestPage_ReadWrite(t *testing.T) {
	t.Parallel()

	const size = 256

	b := raw.NewBuffer(size)

	p := NewPage(size)
	if err := p.Add([]byte{1, 2, 3, 4}); err != nil {
		t.Error(err)
	}

	if err := p.Add([]byte{5, 6, 7, 8}); err != nil {
		t.Error(err)
	}

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
