package pagequeue

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type PageQueue struct {
	pool    Pool
	current *page.Page
	ps      PageStorage
	po      PageOut
}

type PageStorage interface {
	Sync([]byte) error
	Append(record.LogSequenceNumber, []byte) error
	Init([]byte) error
}

func New(
	ps PageStorage,
	size int16,
	headerSize int16,
	slotHeaderSize int16,
) *PageQueue {
	pq := &PageQueue{
		pool: NewPool(size, headerSize, slotHeaderSize),
		ps:   ps,
		po: PageOut{
			ps: ps,
		},
	}

	pq.init()

	return pq
}

func (pq *PageQueue) init() {
	pq.po.pool = &pq.pool
	pq.po.shadow = pq.pool.Get()
	pq.current = pq.pool.Get()
}

func (pq *PageQueue) Init(data []byte) {
	pq.current.Fill(data)
}

func (pq *PageQueue) Count() int {
	return pq.po.Count()
}

// # Append
//   - Compare size of new Record to remaining space in current Page
//   - If it fits, append to the current Page
//   - If it spans, generate a new list of Pages to fit
//   - Append Record as Span to the list
//   - Append list to queue, up to but not including, current
//   - If current Page is entirely full, append current to list and swap current for next Page
func (pq *PageQueue) Append(rc *record.Record) error {
	t := NewTip(&pq.pool, pq.current)
	s := t.Span(int16(rc.Size()))
	if err := rc.Write(s); err != nil {
		return err
	}

	pgs, ok := t.Commit()
	if !ok {
		return ErrUnableToAppend
	}

	pq.appendPages(rc.LogSequenceNumber(), pgs)
	return nil
}

func (pq *PageQueue) appendPages(lsn record.LogSequenceNumber, pgs []*page.Page) {
	pq.po.SetLSN(lsn)

	// If pages is empty, do nothing
	l := len(pgs)
	if l == 0 {
		return
	}

	pq.po.Push(pq.current)
	last := l - 1
	for i, p := range pgs {
		if i == last {
			pq.replaceCurrent(p)
		} else {
			pq.po.Push(p)
		}
	}
}

// # Flushing
//   - Flush entire queue to Disk, including Current
func (pq *PageQueue) Flush() error {
	return pq.po.Flush(pq.current)
}

func (pq *PageQueue) replaceCurrent(p *page.Page) {
	pq.current = p
}

func (pq *PageQueue) Clear() {
	pq.po.Clear()
}
