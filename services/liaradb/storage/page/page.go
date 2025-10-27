package page

import (
	"io"
	"iter"
)

// TODO: Add header data
// - Magic number
// - PageID?
// - TimeLineID
// - Max LogSequenceNumber

// TODO: Potentially use io.OffsetWriter
type Page[H Serializer] struct {
	size   Offset
	header H
	list   List
	items  []Item // TODO: Use Serializer, and use a generic?
}

type Serializer interface {
	Read(io.Reader) error
	Size() int
	Write(io.Writer) error
}

type ReaderAndAt interface {
	io.Reader
	io.ReaderAt
}

type WriterAndAt interface {
	io.Writer
	io.WriterAt
}

type Item = []byte

func New(size Offset) *Page[ZeroHeader] {
	return &Page[ZeroHeader]{
		size:   size,
		header: ZeroHeader{},
	}
}

func NewWithHeader[H Serializer](size Offset, header H) *Page[H] {
	return &Page[H]{
		size:   size,
		header: header,
	}
}

func (p *Page[H]) Add(i Item) error {
	l := len(i)
	if _, err := p.list.Add(p.nextCursor(l), Offset(l)); err != nil {
		// TODO: Test this
		return err
	}

	p.items = append(p.items, i)
	return nil
}

func (p *Page[H]) nextCursor(l int) Offset {
	return Offset(p.Size() - p.list.entriesSize() - l)
}

func (p *Page[H]) Header() H {
	return p.header
}

// TODO: Create a way to iterate rather than reading the entire page
// TODO: Do we need an error parameter?
func (p *Page[H]) Items() iter.Seq2[Item, error] {
	return func(yield func(Item, error) bool) {
		for _, i := range p.items {
			if !yield(i, nil) {
				return
			}
		}
	}
}

func (p *Page[H]) Size() int {
	return int(p.size)
}

// TODO: Should we use seek instead?
func (p *Page[H]) Read(r ReaderAndAt) error {
	if err := p.readHeader(r); err != nil {
		return err
	}

	if err := p.list.Read(r); err != nil {
		return err
	}

	return p.readItems(r)
}

// TODO: Should we use seek instead?
func (p *Page[H]) Write(w WriterAndAt) error {
	if err := p.writeHeader(w); err != nil {
		return err
	}

	if err := p.list.Write(w); err != nil {
		return err
	}

	return p.writeItems(w)
}

func (p *Page[H]) readHeader(r io.Reader) error {
	return p.header.Read(r)
}

func (p *Page[H]) writeHeader(w io.Writer) error {
	return p.header.Write(w)
}

func (p *Page[H]) readItems(r ReaderAndAt) error {
	items := make([]Item, 0, p.list.Length())

	for _, e := range p.list.entries {
		i := make([]byte, e.Length)
		if _, err := r.ReadAt(i, int64(e.Offset)); err != nil {
			return err
		}

		items = append(items, i)
	}

	p.items = items
	return nil
}

func (p *Page[H]) writeItems(w WriterAndAt) error {
	for index, i := range p.items {
		if _, err := w.WriteAt(i, int64(p.list.offset(index))); err != nil {
			return err
		}
	}

	return nil
}
