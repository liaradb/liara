package pagequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type PageQueue struct {
	pool    Pool
	current *page.Page
	po      PageOut
}

func New(
	ps PageStorage,
	size int16,
	headerSize int16,
	slotHeaderSize int16,
) *PageQueue {
	pq := &PageQueue{
		pool: NewPool(size, headerSize, slotHeaderSize),
	}

	pq.init(ps)

	return pq
}

func (pq *PageQueue) init(ps PageStorage) {
	pq.po = newPageOut(ps, pq, &pq.pool)
	pq.current = pq.pool.Get()
}

func (pq *PageQueue) Init(data []byte) {
	pq.current.Fill(data)
}

func (pq *PageQueue) Run(ctx context.Context) {
	pq.po.Run(ctx)
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
	c, ok := pq.po.Append(lsn, pgs...)
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

func (pq *PageQueue) OnFlush(record.LogSequenceNumber) {

}

func (pq *PageQueue) OnError(error) bool {
	return true
}
