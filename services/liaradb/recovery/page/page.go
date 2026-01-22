package page

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

type Page struct {
	page *node
}

func New(size int64) *Page {
	return &Page{
		page: newNode(size),
	}
}

func (p *Page) Init(id action.PageID, tlid action.TimeLineID, rem record.Length) {
	p.page.Reset(id, tlid, rem)
}

func (p *Page) Append(data []byte) bool {
	return p.page.Append(data)
}

func (p *Page) Position() int64 {
	return p.page.Position()
}

func (p *Page) Write(w io.WriterAt) error {
	return p.page.Write(io.NewOffsetWriter(w, p.Position()))
}

func (p *Page) Read(r io.ReadSeeker) error {
	return p.page.Read(r)
}

func (p *Page) ID() action.PageID              { return p.page.ID() }
func (p *Page) TimeLineID() action.TimeLineID  { return p.page.TimeLineID() }
func (p *Page) LengthRemaining() record.Length { return p.page.LengthRemaining() }

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
	return p.page.Read(r)
}

func (p *Page) records() iter.Seq2[*record.Record, error] {
	return func(yield func(*record.Record, error) bool) {
		b := &raw.Buffer{}
		for i := range p.page.Items() {
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
		for i := range p.page.ItemsReverse() {
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
