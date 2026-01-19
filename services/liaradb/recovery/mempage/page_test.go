package mempage

import (
	"io"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/raw"
)

func TestPage_Add(t *testing.T) {
	t.Parallel()

	const size = 256
	p := New(size)

	items := []*item{
		newItem([]byte{1, 2, 3, 4}),
		newItem([]byte{5, 6, 7, 8})}
	for _, i := range items {
		if _, err := p.Add(i.data); err != nil {
			t.Error(err)
		}
	}

	result := make([]*item, 0)

	for i, err := range p.Items() {
		if err != nil {
			t.Error(err)
		}

		result = append(result, newItem(i))
	}

	if !slices.EqualFunc(result, items, func(a, b *item) bool {
		return a.Compare(b)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", result, items)
	}
}

func TestPage_Add__ErrInsufficientSpace(t *testing.T) {
	t.Parallel()

	const size = 16
	p := New(size)

	items := []*item{
		newItem([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
		newItem([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17})}

	for _, i := range items {
		if _, err := p.Add(i.data); err != raw.ErrInsufficientSpace {
			t.Error("should return insufficient space")
		}
	}
}

func TestPage_ReadWrite(t *testing.T) {
	t.Parallel()

	const size = 256
	p := New(size)

	b, items := writeRecords(t, size, p)

	p1 := New(size)

	if err := p1.Read(b); err != nil {
		t.Fatal(err)
	}

	result := make([]*item, 0)

	for i, err := range p1.Items() {
		if err != nil {
			t.Error(err)
		}

		result = append(result, newItem(i))
	}

	if !slices.EqualFunc(result, items, func(a, b *item) bool {
		return a.Compare(b)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", result, items)
	}
}

func TestPage_Reverse(t *testing.T) {
	t.Parallel()

	const size = 256
	p := New(size)

	b, items := writeRecords(t, size, p)

	p1 := New(size)

	if err := p1.Read(b); err != nil {
		t.Fatal(err)
	}

	result := make([]*item, 0)

	for i, err := range p1.ItemsReverse() {
		if err != nil {
			t.Error(err)
		}

		result = append(result, newItem(i))
	}

	slices.Reverse(items)

	if !slices.EqualFunc(result, items, func(a, b *item) bool {
		return a.Compare(b)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", result, items)
	}
}

// TODO: Test this
func TestPage_ReadWrite__Header(t *testing.T) {
	t.Skip()
	t.Parallel()

	// const size = 256
	// data := raw.BaseString("test data")
	// p := New(size, &testPageHeader{data})

	// b, items := writeRecords(t, size, p)

	// p1 := New(size, &testPageHeader{})

	// if err := p1.Read(b); err != nil {
	// 	t.Fatal(err)
	// }

	// result := make([]*item, 0)

	// for i, err := range p1.Items() {
	// 	if err != nil {
	// 		t.Error(err)
	// 	}

	// 	result = append(result, newItem(i))
	// }

	// if !slices.EqualFunc(result, items, func(a, b *item) bool {
	// 	return a.Compare(b)
	// }) {
	// 	t.Errorf("incorrect result: %v, expected: %v", result, items)
	// }

	// if d := p.Header().data; d != data {
	// 	t.Errorf("incorrect header: %v, expected: %v", d, data)
	// }
}

func writeRecords(t *testing.T, size int64, p *Page) (*raw.Buffer, []*item) {
	b := raw.NewBuffer(size)
	items := createRecords(4, 32)

	for _, i := range items {
		if _, err := p.Add(i.data); err != nil {
			t.Fatal(err)
		}
	}

	b.Clear()
	if err := p.Write(b); err != nil {
		t.Fatal(err)
	}

	if s := p.Size(); int64(s) != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	if _, err := b.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	return b, items
}

func createRecords(rows, count int) []*item {
	items := make([]*item, 0, rows)
	for i := range byte(cap(items)) {
		item := make([]byte, 0, count)
		for j := range byte(cap(item)) {
			item = append(item, j+i*byte(cap(item)))
		}
		items = append(items, newItem(item))
	}
	return items
}

type testPageHeader struct {
	data raw.BaseString
}

var _ Serializer = (*testPageHeader)(nil)

func (t *testPageHeader) Read(r io.Reader) error {
	return t.data.Read(r)
}

func (t *testPageHeader) Size() int {
	return t.data.Size()
}

func (t *testPageHeader) Write(w io.Writer) error {
	return t.data.Write(w)
}
