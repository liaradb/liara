package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/crclist"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

const (
	itemSize = 8
)

type Page struct {
	header
	data     []byte
	list     crclist.CRCList
	byteList bytelist.ByteList
}

func New(size int64) *Page {
	data := make([]byte, size)
	return NewFromSlice(data)
}

func NewFromSlice(data []byte) *Page {
	header, data0 := newHeader(data)
	return &Page{
		header:   header,
		data:     data,
		list:     crclist.New(data0),
		byteList: bytelist.New(data0),
	}
}

// TODO: Test this
func (p *Page) Init(pid action.PageID, tlid action.TimeLineID, rl record.Length) {
	clear(p.data)
	p.list.Clear()
	p.header.Reset(pid, tlid, rl)
}

func (p *Page) Append(data []byte) bool {
	size := int16(len(data))
	if !p.hasSpace(size) {
		return false
	}

	offset := p.next() - size
	if _, ok := p.list.Push(offset, size, page.NewCRC(data)); !ok {
		return false
	}

	p.header.setNext(offset)

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok { // We already checked hasSpace
		return false
	}

	copy(b, data)

	return true
}

func (p *Page) Space() int16 {
	next := p.next()
	size := p.list.Size()
	return max(next-size-itemSize, 0)
}

func (p *Page) next() int16 {
	size := p.list.Count()
	if size == 0 {
		return int16(p.list.Length())
	} else {
		return p.header.Next()
	}
}

func (p *Page) hasSpace(size int16) bool {
	s := p.Space()
	return size <= s
}

func (p *Page) Position() int64 {
	return p.header.ID().Position(int64(len(p.data)))
}

func (p *Page) Write(w io.WriterAt) error {
	wr := io.NewOffsetWriter(w, p.Position())
	_, err := wr.Write(p.data)
	return err
}

func (p *Page) Read(r io.ReadSeeker) error {
	_, err := r.Read(p.data)
	if err != nil {
		return err
	}

	p.list.Clear()
	return nil
}

func (p *Page) Iterate(r io.ReadSeeker) (iter.Seq2[*record.Record, error], error) {
	if err := p.read(r); err != nil {
		return nil, err
	}

	return p.records(), nil
}

func (p *Page) Reverse(r io.ReadSeeker) (iter.Seq2[*record.Record, error], error) {
	if err := p.read(r); err != nil {
		return nil, err
	}

	return p.reverse(), nil
}

func (p *Page) read(r io.ReadSeeker) error {
	return p.Read(r)
}

func (p *Page) records() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		b := &raw.Buffer{}
		for i := range p.Items() {
			b.Reset(i)
			rc := &record.Record{}
			if err := rc.Read(b); err != nil {
				yield(nil, err)
				return
			}

			if !yield(rc, nil) {
				return
			}
		}
	}
}

func (p *Page) reverse() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		b := &raw.Buffer{}
		for i := range p.ItemsReverse() {
			b.Reset(i)
			rc := &record.Record{}
			if err := rc.Read(b); err != nil {
				yield(nil, err)
				return
			}

			if !yield(rc, nil) {
				return
			}
		}
	}
}

func (p *Page) Child(index int16) ([]byte, bool) {
	i, ok := p.list.Item(index)
	if !ok {
		return nil, false
	}

	d, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
	if !ok {
		return nil, false
	}

	if !i.CRC.Compare(d) {
		return nil, false
	}

	return d, true
}

func (p *Page) Items() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.Items() {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}

func (p *Page) ItemsReverse() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.ItemsReverse() {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}

func (p *Page) ChildrenRange(start, end int16) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range p.list.ItemsRange(start, end) {
			b, ok := p.byteList.Slice(int64(i.Offset), int64(i.Size))
			if !ok || !i.CRC.Compare(b) || !yield(b) {
				return
			}
		}
	}
}
