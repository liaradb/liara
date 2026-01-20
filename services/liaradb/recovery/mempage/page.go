package mempage

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

// TODO: Add header data
// - Magic number
// - PageID?
// - TimeLineID
// - Max LogSequenceNumber

// TODO: Use List to return slices of the underlying slice
// Don't parse items here, just return slices
// Then use the slices in a [raw.Buffer] to allow reading.

// TODO: Potentially use io.OffsetWriter
type Page struct {
	size   page.Offset
	header *header
	list   list
	items  []*item // TODO: Change back to []byte
}

type Serializer interface {
	Read(io.Reader) error
	Size() int
	Write(io.Writer) error
}

func New(size int64) *Page {
	header := &header{}
	return &Page{
		size:   page.Offset(size),
		header: header,
		list:   newList(page.MagicSize + header.Size()),
	}
}

func (p *Page) ID() action.PageID              { return p.header.ID() }
func (p *Page) TimeLineID() action.TimeLineID  { return p.header.TimeLineID() }
func (p *Page) LengthRemaining() record.Length { return p.header.LengthRemaining() }

func (p *Page) Reset(
	id action.PageID,
	timeLineID action.TimeLineID,
	lengthRemaining record.Length,
) {
	p.header = newHeader(id, timeLineID, lengthRemaining)
	p.list.Reset()
	p.items = nil
}

// TODO: Test offset return
func (p *Page) Add(data []byte) (page.Offset, error) {
	i := newItem(data)
	l := i.Size()
	offset, err := p.list.Add(p.nextCursor(l), listLength(l))
	if err != nil {
		// TODO: Test this
		return 0, err
	}

	p.items = append(p.items, i)
	return offset, nil
}

func (p *Page) nextCursor(l int) page.Offset {
	return page.Offset(p.Size() - p.list.entriesSize() - l)
}

func (p *Page) Size() int { return int(p.size) }

// TODO: Create a way to iterate rather than reading the entire page
func (p *Page) Items() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for _, i := range p.items {
			if !yield(i.data) {
				return
			}
		}
	}
}

func (p *Page) ItemsReverse() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		l := len(p.items) - 1
		for index := range p.items {
			if !yield(p.items[l-index].data) {
				return
			}
		}
	}
}

func (p *Page) Read(r io.ReadSeeker) error {
	if err := p.readHeader(r); err != nil {
		return err
	}

	if err := p.list.Read(r); err != nil {
		return err
	}

	return p.readItems(r)
}

func (p *Page) Write(w io.WriteSeeker) error {
	if err := p.writeItems(w); err != nil {
		return err
	}

	if _, err := w.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := p.writeHeader(w); err != nil {
		return err
	}

	return p.list.Write(w)
}

func (p *Page) readHeader(r io.Reader) error {
	var m page.Magic
	return raw.ReadAll(r,
		&m,
		p.header)
}

func (p *Page) writeHeader(w io.Writer) error {
	return raw.WriteAll(w,
		page.MagicPage,
		p.header)
}

func (p *Page) readItems(r io.ReadSeeker) error {
	items := make([]*item, 0, p.list.Length())

	for _, e := range p.list.entries {
		if _, err := r.Seek(int64(e.Offset), io.SeekStart); err != nil {
			return err
		}

		i := newItemByLength(e.Length)
		if err := i.Read(r, e.CRC); err != nil {
			return err
		}

		items = append(items, i)
	}

	p.items = items
	return nil
}

func (p *Page) writeItems(w io.WriteSeeker) error {
	for index, i := range p.items {
		if _, err := w.Seek(int64(p.list.offset(index)), io.SeekStart); err != nil {
			return err
		}

		// TODO: Don't calculate CRC if we have it already
		crc, err := i.Write(w)
		if err != nil {
			return err
		}

		p.list.setCRC(index, crc)
	}

	return nil
}
