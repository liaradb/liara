package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
)

type WriteQueue struct {
	items chan queueItem
	ps    PageStorage
	pool  *pagepool.PagePool
}

type queueItem interface {
	Wait(context.Context) error
	Store(PageStorage) error
	Page() *page.Page
}

type PageStorage interface {
	Replace([]byte) error
	Append(record.LogSequenceNumber, []byte) error
	Init([]byte) error
}

func New(
	size int,
	ps PageStorage,
	pool *pagepool.PagePool,
) *WriteQueue {
	return &WriteQueue{
		items: make(chan queueItem, size),
		ps:    ps,
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
	if err := qi.Store(wq.ps); err != nil {
		return err
	}

	wq.pool.Put(qi.Page())
	return nil
}

func (wq *WriteQueue) Append(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	p *page.Page,
) {
	select {
	case wq.items <- newAppendQueueItem(lsn, p):
	case <-ctx.Done():
	}
}

func (wq *WriteQueue) AppendSync(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	p *page.Page,
) error {
	qi := newAppendSyncQueueItem(lsn, p)

	select {
	case wq.items <- qi:
	case <-ctx.Done():
	}

	return qi.Wait(ctx)
}

func (wq *WriteQueue) Sync(
	ctx context.Context,
	p *page.Page,
) error {
	qi := newSyncQueueItem(p)

	select {
	case wq.items <- qi:
	case <-ctx.Done():
	}

	return qi.Wait(ctx)
}
