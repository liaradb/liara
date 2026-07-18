package pagequeue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/pagequeue/writequeue"
	"github.com/liaradb/liaradb/recovery/record"
)

type PageQueue struct {
	pool    *pagepool.PagePool
	wq      *writequeue.WriteQueue
	ps      writequeue.PageStorage
	current *logpage.LogPage
	flushed bool
	lsn     record.LogSequenceNumber
}

func New(
	ps writequeue.PageStorage,
	pl *pagepool.PagePool,
	writeQueueSize int,
) *PageQueue {
	return &PageQueue{
		pool:    pl,
		wq:      writequeue.New(writeQueueSize, ps, pl),
		ps:      ps,
		current: pl.Get(),
	}
}

func (pq *PageQueue) Init(data []byte) error {
	pq.current.Fill(data)
	pq.flushed = true
	return nil
}

// TODO: Test this error
func (pq *PageQueue) Run(ctx context.Context) error {
	return pq.wq.Run(ctx)
}

// # Append
//   - Compare size of new Record to remaining space in current Page
//   - If it fits, append to the current Page
//   - If it spans, generate a new list of Pages to fit
//   - Append Record as Span to the list
//   - Append list to queue, up to but not including, current
//   - If current Page is entirely full, append current to list and swap current for next Page
func (pq *PageQueue) Append(ctx context.Context, rc *record.Record) error {
	return pq.AppendWait(ctx, rc, nil)
}

func (pq *PageQueue) AppendWait(ctx context.Context, rc *record.Record, h func()) error {
	t := NewTip(pq.pool, pq.current)
	s := t.Span(rc.Size())
	if err := rc.Write(s); err != nil {
		return err
	}

	s.Commit()
	pgs, ok := t.Commit(rc.LogSequenceNumber(), h)
	if !ok {
		return writequeue.ErrUnableToAppend
	}

	pq.appendPages(ctx, rc.LogSequenceNumber(), pgs)
	return nil
}

func (pq *PageQueue) appendPages(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	pgs []*logpage.LogPage,
) {
	// If only one page, do nothing
	l := len(pgs)
	if l <= 1 {
		return
	}

	pq.flushCurrent(ctx, lsn, pgs[0])
	for _, p := range pgs[1 : l-1] {
		pq.wq.Append(ctx, lsn, p)
	}

	pq.replaceCurrent(lsn, pgs[l-1])
}

func (pq *PageQueue) flushCurrent(ctx context.Context, lsn record.LogSequenceNumber, current *logpage.LogPage) {
	if pq.flushed {
		pq.wq.Replace(ctx, current)
	} else {
		pq.wq.Append(ctx, lsn, current)
		pq.flushed = true
	}
}

func (pq *PageQueue) replaceCurrent(lsn record.LogSequenceNumber, p *logpage.LogPage) {
	pq.lsn = lsn
	pq.current = p
	pq.flushed = false
}

// # Flushing
//   - Flush entire queue to Disk, including Current
func (pq *PageQueue) Flush(ctx context.Context) error {
	shadow := pq.pool.Get()
	shadow.Shadow(pq.current)
	if pq.flushed {
		if err := pq.wq.ReplaceSync(ctx, shadow); err != nil {
			return err
		}
	} else {
		if err := pq.wq.AppendSync(ctx, pq.lsn, shadow); err != nil {
			return err
		}

		pq.flushed = true
	}

	return nil
}
