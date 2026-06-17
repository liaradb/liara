package pagequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type WriteQueue struct {
	items chan queueItem
	ps    PageStorage
}

type queueItem interface {
	Wait(context.Context) error
	Store(PageStorage) error
}

type PageStorage interface {
	Sync([]byte) error
	Append(record.LogSequenceNumber, []byte) error
	Init([]byte) error
}

func newWriteQueue(
	size int,
	ps PageStorage,
) *WriteQueue {
	return &WriteQueue{
		items: make(chan queueItem, size),
		ps:    ps,
	}
}

func (wq *WriteQueue) Run(ctx context.Context) {
	go wq.run(ctx)
}

func (wq *WriteQueue) run(ctx context.Context) error {
	for {
		select {
		case qi := <-wq.items:
			if err := wq.runItem(qi); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (wq *WriteQueue) runItem(qi queueItem) error {
	return qi.Store(wq.ps)
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

func (wq *WriteQueue) Sync(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	p *page.Page,
) error {
	qi := newSyncQueueItem(p)

	select {
	case wq.items <- qi:
	case <-ctx.Done():
	}

	return qi.Wait(ctx)
}
