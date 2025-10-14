package record

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type Page struct {
	size   Offset
	list   List
	items  []Item
	buffer *raw.Buffer
}

type Item = []byte

func NewPage(size Offset) *Page {
	b := raw.NewBuffer(int(size))
	return &Page{
		size:   size,
		buffer: b,
	}
}

func (p *Page) Add(i Item) {
	l := len(i)
	cursor := p.Size() - p.list.entriesSize() - l
	p.list.Add(Offset(cursor), Offset(l))
	p.items = append(p.items, i)
}

func (p *Page) Size() int {
	return int(p.size)
}

func (p *Page) Write(w interface {
	// io.WriterAt
	io.Writer
}) error {
	p.buffer.Clear()

	if err := p.list.Write(p.buffer); err != nil {
		return err
	}

	for index, i := range p.items {
		if err := p.writeItem(p.buffer, i, p.list.offset(index)); err != nil {
			return err
		}
	}

	return p.writeBuffer(w)
}

func (p *Page) writeItem(w io.WriterAt, item Item, off Offset) error {
	_, err := w.WriteAt(item, int64(off))
	return err
}

func (p *Page) writeBuffer(w io.Writer) error {
	if _, err := p.buffer.Seek(0, io.SeekStart); err != nil {
		return err
	}

	_, err := io.Copy(w, p.buffer)
	p.buffer.Clear()
	return err
}

func (p *Page) Read(r io.Reader) error {
	p.buffer.Clear()

	if err := p.readBuffer(r); err != nil {
		return err
	}

	if err := p.list.Read(p.buffer); err != nil {
		return err
	}

	for _, e := range p.list.entries {
		i := make([]byte, e.Length)
		if _, err := p.buffer.ReadAt(i, int64(e.Offset)); err != nil {
			return err
		}

		p.items = append(p.items, i)
	}

	p.buffer.Clear()
	return nil
}

func (p *Page) readBuffer(r io.Reader) error {
	if _, err := io.Copy(p.buffer, r); err != nil {
		return err
	}

	_, err := p.buffer.Seek(0, io.SeekStart)
	return err
}
