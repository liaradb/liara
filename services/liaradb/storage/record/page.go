package record

import (
	"io"
	"iter"
)

type Page struct {
	size  Offset
	list  List
	items []Item
}

type Item = []byte

func NewPage(size Offset) *Page {
	return &Page{
		size: size,
	}
}

func (p *Page) Add(i Item) error {
	l := len(i)
	if _, err := p.list.Add(p.nextCursor(l), Offset(l)); err != nil {
		return err
	}

	p.items = append(p.items, i)
	return nil
}

func (p *Page) nextCursor(l int) Offset {
	return Offset(p.Size() - p.list.entriesSize() - l)
}

// TODO: Do we need an error parameter?
func (p *Page) Items() iter.Seq2[Item, error] {
	return func(yield func(Item, error) bool) {
		for _, i := range p.items {
			if !yield(i, nil) {
				return
			}
		}
	}
}

func (p *Page) Size() int {
	return int(p.size)
}

// TODO: Should we use seek instead?
func (p *Page) Write(w interface {
	io.WriterAt
	io.Writer
}) error {
	if err := p.list.Write(w); err != nil {
		return err
	}

	for index, i := range p.items {
		if _, err := w.WriteAt(i, int64(p.list.offset(index))); err != nil {
			return err
		}
	}

	return nil
}

// TODO: Should we use seek instead?
func (p *Page) Read(r interface {
	io.Reader
	io.ReaderAt
}) error {
	if err := p.list.Read(r); err != nil {
		return err
	}

	p.items = make([]Item, 0, p.list.Length())

	for _, e := range p.list.entries {
		i := make([]byte, e.Length)
		if _, err := r.ReadAt(i, int64(e.Offset)); err != nil {
			return err
		}

		p.items = append(p.items, i)
	}

	return nil
}
