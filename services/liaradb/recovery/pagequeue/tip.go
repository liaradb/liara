package pagequeue

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type Tip struct {
	size           int16
	headerSize     int16
	slotHeaderSize int16
	current        *page.Page
	pages          []*page.Page
	sizes          []int16
}

func NewTip(size int16, headerSize int16, slotHeaderSize int16, current *page.Page) Tip {
	return Tip{
		size:           size,
		headerSize:     headerSize,
		slotHeaderSize: slotHeaderSize,
		current:        current,
	}
}

func (t *Tip) Pages() []*page.Page { return t.pages }

// Request Lease from current Page
// If insufficient space is available, build list of Pages for remaining
func (t *Tip) Span(size int16) *record.Span {
	s := record.NewSpan()

	var available int16 = 0
	var remaining int16 = size

	_, data := t.current.Next(remaining)

	l := int16(len(data))
	t.sizes = append(t.sizes, l)
	available = l
	remaining -= l
	f := record.NewFragment(data)
	s.Append(f)

	for available < size {
		p := t.next()
		_, data := p.Next(remaining)

		l := int16(len(data))
		t.sizes = append(t.sizes, l)
		available += l
		remaining -= l
		f := record.NewFragment(data)
		s.Append(f)
	}

	return &s
}

func (t *Tip) next() *page.Page {
	p := page.New(t.size, t.headerSize, t.slotHeaderSize)
	t.pages = append(t.pages, p)
	return p
}

// TODO: What do we do on a partial commit?
func (t *Tip) Commit() bool {
	size := t.sizes[0]
	if _, ok := t.current.Commit(size); !ok {
		return false
	}

	for i, p := range t.pages {
		size := t.sizes[i]
		if _, ok := p.Commit(size); !ok {
			return false
		}
	}

	return true
}
