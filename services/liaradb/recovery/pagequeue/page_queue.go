package pagequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/span"
	"github.com/liaradb/liaradb/recovery/writequeue"
)

type PageQueue struct {
	pool    pagepool.PagePool
	current *page.Page
	po      PageOut
}

func New(
	ps PageStorage,
	size int16,
	headerSize int16,
) *PageQueue {
	pq := &PageQueue{
		pool: pagepool.New(size, headerSize, span.FragmentHeaderSize),
	}

	pq.init(ps)

	return pq
}

func (pq *PageQueue) init(ps PageStorage) {
	pq.po = newPageOut(ps, &pq.pool)
	pq.current = pq.pool.Get()
}

func (pq *PageQueue) Init(data []byte) {
	pq.current.Fill(data)
}

// TODO: Test this error
func (pq *PageQueue) Run(ctx context.Context) error {
	return pq.po.Run(ctx)
}

// # Append
//   - Compare size of new Record to remaining space in current Page
//   - If it fits, append to the current Page
//   - If it spans, generate a new list of Pages to fit
//   - Append Record as Span to the list
//   - Append list to queue, up to but not including, current
//   - If current Page is entirely full, append current to list and swap current for next Page
func (pq *PageQueue) Append(ctx context.Context, rc *record.Record) error {
	t := NewTip(&pq.pool, pq.current)
	s := t.Span(int16(rc.Size()))
	if err := rc.Write(s); err != nil {
		return err
	}

	pgs, ok := t.Commit()
	if !ok {
		return writequeue.ErrUnableToAppend
	}

	pq.appendPages(ctx, rc.LogSequenceNumber(), pgs)
	return nil
}

func (pq *PageQueue) appendPages(ctx context.Context, lsn record.LogSequenceNumber, pgs []*page.Page) {
	c, ok := pq.po.Append(ctx, lsn, pgs...)
	if ok {
		pq.replaceCurrent(c)
	}
}

func (pq *PageQueue) replaceCurrent(p *page.Page) {
	pq.current = p
}

func (pq *PageQueue) Clear() {
	pq.po.Clear()
}

func (pq *PageQueue) Count() int {
	return pq.po.Count()
}

// # Flushing
//   - Flush entire queue to Disk, including Current
func (pq *PageQueue) Flush(ctx context.Context) error {
	return pq.po.Flush(ctx, pq.current)
}
