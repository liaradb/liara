package record

import "io"

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

func (p *Page) Add(i Item) {
	l := len(i)
	p.list.Add(p.nextCursor(l), Offset(l))
	p.items = append(p.items, i)
}

func (p *Page) nextCursor(l int) Offset {
	return Offset(p.Size() - p.list.entriesSize() - l)
}

func (p *Page) Size() int {
	return int(p.size)
}

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

func (p *Page) Read(r interface {
	io.Reader
	io.ReaderAt
}) error {
	if err := p.list.Read(r); err != nil {
		return err
	}

	for _, e := range p.list.entries {
		i := make([]byte, e.Length)
		if _, err := r.ReadAt(i, int64(e.Offset)); err != nil {
			return err
		}

		p.items = append(p.items, i)
	}

	return nil
}
