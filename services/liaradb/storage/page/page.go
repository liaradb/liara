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
type Page[H Serializer, I Serializer] struct {
	size   Offset
	header H
	list   List
	items  []I
	newI   func(Offset) I
}

type Serializer interface {
	Read(io.Reader) error
	Size() int
	Write(io.Writer) error
}

type Item = []byte

func New(size Offset) *Page[ZeroHeader, *ItemSerializer] {
	return &Page[ZeroHeader, *ItemSerializer]{
		size:   size,
		header: ZeroHeader{},
		newI:   NewItemSerializerByLength,
	}
}

func NewWithHeader[H Serializer, I Serializer](size Offset, header H, newI func(Offset) I) *Page[H, I] {
	return &Page[H, I]{
		size:   size,
		header: header,
		newI:   newI,
	}
}

func (p *Page[H, I]) Add(i I) error {
	l := i.Size()
	if _, err := p.list.Add(p.nextCursor(l), Offset(l)); err != nil {
		// TODO: Test this
		return err
	}

	p.items = append(p.items, i)
	return nil
}

func (p *Page[H, I]) nextCursor(l int) Offset {
	return Offset(p.Size() - p.list.entriesSize() - l)
}

func (p *Page[H, I]) Header() H {
	return p.header
}

// TODO: Create a way to iterate rather than reading the entire page
// TODO: Do we need an error parameter?
func (p *Page[H, I]) Items() iter.Seq2[I, error] {
	return func(yield func(I, error) bool) {
		for _, i := range p.items {
			if !yield(i, nil) {
				return
			}
		}
	}
}

func (p *Page[H, I]) Size() int {
	return int(p.size)
}

// TODO: Should we use seek instead?
func (p *Page[H, I]) Read(r io.ReadSeeker) error {
	if err := p.readHeader(r); err != nil {
		return err
	}

	if err := p.list.Read(r); err != nil {
		return err
	}

	return p.readItems(r)
}

// TODO: Should we use seek instead?
func (p *Page[H, I]) Write(w io.WriteSeeker) error {
	if err := p.writeHeader(w); err != nil {
		return err
	}

	if err := p.list.Write(w); err != nil {
		return err
	}

	return p.writeItems(w)
}

func (p *Page[H, I]) readHeader(r io.Reader) error {
	return p.header.Read(r)
}

func (p *Page[H, I]) writeHeader(w io.Writer) error {
	return p.header.Write(w)
}

func (p *Page[H, I]) readItems(r io.ReadSeeker) error {
	items := make([]I, 0, p.list.Length())

	for _, e := range p.list.entries {
		if _, err := r.Seek(int64(e.Offset), io.SeekStart); err != nil {
			return err
		}

		i := p.newI(e.Length)
		if err := i.Read(r); err != nil {
			return err
		}

		items = append(items, i)
	}

	p.items = items
	return nil
}

func (p *Page[H, I]) writeItems(w io.WriteSeeker) error {
	for index, i := range p.items {
		if _, err := w.Seek(int64(p.list.offset(index)), io.SeekStart); err != nil {
			return err
		}

		if err := i.Write(w); err != nil {
			return err
		}
	}

	return nil
}
