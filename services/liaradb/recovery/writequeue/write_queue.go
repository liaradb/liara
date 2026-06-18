package writequeue

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

func New(
	size int,
	ps PageStorage,
) *WriteQueue {
	return &WriteQueue{
		items: make(chan queueItem, size),
		ps:    ps,
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

// TODO: Return Page to Pool after it is stored
func (wq *WriteQueue) run(qi queueItem) error {
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
