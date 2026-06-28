package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
)

type WriteQueue struct {
	items chan queueItem
	ps    PageStorage
	fl    Flusher
	pool  *pagepool.PagePool
}

type queueItem interface {
	Wait(context.Context) error
	Store(PageStorage, Flusher) error
	Page() *logpage.LogPage
}

type PageStorage interface {
	Replace([]byte) error
	Append(record.LogSequenceNumber, []byte) error
	Init([]byte) error
}

type Flusher interface {
	OnFlush(record.LogSequenceNumber)
}

func New(
	size int,
	ps PageStorage,
	fl Flusher,
	pool *pagepool.PagePool,
) *WriteQueue {
	return &WriteQueue{
		items: make(chan queueItem, size),
		ps:    ps,
		fl:    fl,
		pool:  pool,
	}
}

func (wq *WriteQueue) Run(ctx context.Context) error {
	for {
		select {
		case qi := <-wq.items:
			if err := wq.run(qi); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (wq *WriteQueue) run(qi queueItem) error {
	if err := qi.Store(wq.ps, wq.fl); err != nil {
		return err
	}

	wq.pool.Put(qi.Page())
	return nil
}

func (wq *WriteQueue) Append(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	p *logpage.LogPage,
) {
	select {
	case wq.items <- newAppendQueueItem(lsn, p):
	case <-ctx.Done():
	}
}

func (wq *WriteQueue) AppendSync(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	p *logpage.LogPage,
) error {
	qi := newAppendSyncQueueItem(lsn, p)

	select {
	case wq.items <- qi:
	case <-ctx.Done():
	}

	return qi.Wait(ctx)
}

func (wq *WriteQueue) Replace(
	ctx context.Context,
	p *logpage.LogPage,
) {
	select {
	case wq.items <- newReplaceQueueItem(p):
	case <-ctx.Done():
	}
}

func (wq *WriteQueue) ReplaceSync(
	ctx context.Context,
	p *logpage.LogPage,
) error {
	qi := newReplaceSyncQueueItem(p)

	select {
	case wq.items <- qi:
	case <-ctx.Done():
	}

	return qi.Wait(ctx)
}
