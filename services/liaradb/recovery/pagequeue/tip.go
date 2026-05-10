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
	p := t.pool.Get()
	t.pages = append(t.pages, p)
	return p
}

// Commit pages before current to avoid a partial commit
func (t *Tip) Commit() bool {
	for i, p := range t.pages {
		size := t.sizes[i]
		if _, ok := p.Commit(size); !ok {
			return false
		}
	}

	size := t.sizes[0]
	if _, ok := t.current.Commit(size); !ok {
		return false
	}

	return true
}
