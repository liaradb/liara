package page

import (
	"io"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/raw"
)

func TestPage_Add(t *testing.T) {
	t.Parallel()

	const size = 256
	p := New(size)

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

	p := New(size)

	items := createRecords(2, 16)
	for _, i := range items {
		if err := p.Add(i); err != nil {
			t.Error(err)
		}
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

	p1 := New(256)

	if err := p1.Read(b); err != nil {
		t.Fatal(err)
	}

	result := make([]Item, 0)

	for i, err := range p1.Items() {
		if err != nil {
			t.Error(err)
		}

		result = append(result, i)
	}

	if !slices.EqualFunc(result, items, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, items)
	}
}

func createRecords(rows, count int) [][]byte {
	items := make([][]byte, 0, rows)
	for i := range byte(cap(items)) {
		item := make([]byte, 0, count)
		for j := range byte(cap(item)) {
			item = append(item, j+i*byte(cap(item)))
		}
		items = append(items, item)
	}
	return items
}
