package pagequeue

import (
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type Tip struct {
	current *page.Page
	pages   []*page.Page
}

func NewTip(current *page.Page) Tip {
	return Tip{
		current: current,
		pages:   []*page.Page{current},
	}
}

// Request Lease from current Page
// If insufficient space is available, build list of Pages for remaining
func (t *Tip) Span(size int16) *record.Span {
	s := record.NewSpan()

	var available int16 = 0
	var remaining int16 = size

	data, ok := t.current.Lease(remaining)
	if !ok {
		// TODO: This ok may not be correct
		panic("incorrect")
	}

	l := int16(len(data))
	available = l
	remaining -= l
	f := record.NewFragment(data)
	s.Append(f)

	for available < size {
		p := t.next()
		data, ok := p.Lease(remaining)
		if !ok {
			// TODO: This ok may not be correct
			panic("incorrect")
		}

		l := int16(len(data))
		available += l
		remaining -= l
		f := record.NewFragment(data)
		s.Append(f)
	}

	return &s
}

func (t *Tip) next() *page.Page {
	p := page.New(int64(t.current.Size()))
	p.Init(t.current.ID()+1, t.current.TimeLineID(), record.NewLength(0))
	t.pages = append(t.pages, p)
	return p
}

func (t *Tip) Pages() []*page.Page {
	return t.pages
}
