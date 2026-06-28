package pagequeue

import (
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/span"
)

type Tip struct {
	pool    *pagepool.PagePool
	current *logpage.LogPage
	pages   []*logpage.LogPage
	sizes   []int16
}

func NewTip(pool *pagepool.PagePool, current *logpage.LogPage) Tip {
	return Tip{
		pool:    pool,
		current: current,
	}
}

// Request Lease from current Page
// If insufficient space is available, build list of Pages for remaining
func (t *Tip) Span(size int16) *span.Span {
	s := span.Span{}

	t.pages = append(t.pages, t.current)

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

	s.InitIndexes()
	return &s
}

func (t *Tip) appendToSpan(s *span.Span, p *logpage.LogPage, remaining int16) int16 {
	header, data := p.Next(remaining)
	l := int16(len(data))
	t.sizes = append(t.sizes, l)
	if l == 0 {
		return l
	}

	_ = s.Append(header, data)
	return l
}

func (t *Tip) next() *logpage.LogPage {
	p := t.pool.Get()
	t.pages = append(t.pages, p)
	return p
}

func (t *Tip) Commit(lsn record.LogSequenceNumber) ([]*logpage.LogPage, bool) {
	if ok := t.commitPages(); !ok {
		t.abortPages()
		return nil, false
	}

	t.pages[len(t.pages)-1].SetLSN(lsn)
	return t.pages, true
}

// Commit pages before current to avoid a partial commit
func (t *Tip) commitPages() bool {
	for i, p := range t.pages[1:] {
		if !t.commitPage(p, i+1) {
			return false
		}
	}

	return t.commitPage(t.current, 0)
}

func (t *Tip) commitPage(p *logpage.LogPage, i int) bool {
	size := t.sizes[i]
	if size == 0 {
		return true
	}

	return p.Commit(size)
}

// Put everything after the first page back
func (t *Tip) abortPages() {
	for _, p := range t.pages[1:] {
		t.pool.Put(p)
	}
}
