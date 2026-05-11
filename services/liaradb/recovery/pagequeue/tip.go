package pagequeue

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type Tip struct {
	pool    *Pool
	current *page.Page
	pages   []*page.Page
	sizes   []int16
}

func NewTip(pool *Pool, current *page.Page) Tip {
	return Tip{
		pool:    pool,
		current: current,
	}
}

// Request Lease from current Page
// If insufficient space is available, build list of Pages for remaining
func (t *Tip) Span(size int16) *record.Span {
	s := record.NewSpan()

	var available int16 = 0
	var remaining int16 = size

	p := t.current

	l := t.appendToSpan(&s, p, remaining)

	available = l
	remaining -= l

	for available < size {
		p = t.next()

		l := t.appendToSpan(&s, p, remaining)

		available += l
		remaining -= l
	}

	return &s
}

func (t *Tip) appendToSpan(s *record.Span, p *page.Page, remaining int16) int16 {
	_, data := p.Next(remaining)
	l := int16(len(data))
	t.sizes = append(t.sizes, l)
	f := record.NewFragment(data)
	s.Append(f)
	return l
}

func (t *Tip) next() *page.Page {
	p := t.pool.Get()
	t.pages = append(t.pages, p)
	return p
}

// Commit pages before current to avoid a partial commit
func (t *Tip) Commit() ([]*page.Page, bool) {
	for i, p := range t.pages {
		size := t.sizes[i]
		if _, ok := p.Commit(size); !ok {
			return nil, false
		}
	}

	size := t.sizes[0]
	if _, ok := t.current.Commit(size); !ok {
		return nil, false
	}

	return t.pages, true
}

func (t *Tip) Abort() {
	for _, p := range t.pages {
		t.pool.Put(p)
	}
}
