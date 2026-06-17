package pagequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type WriteQueue struct {
	items chan queueItem
	ps    PageStorage
	fl    Flusher
}

type PageStorage interface {
	Sync([]byte) error
	Append(record.LogSequenceNumber, []byte) error
	Init([]byte) error
}

type Flusher interface {
	OnFlush(record.LogSequenceNumber)
	OnError(error) bool
}

func newWriteQueue(
	size int,
	ps PageStorage,
	fl Flusher,
) *WriteQueue {
	return &WriteQueue{
		items: make(chan queueItem, size),
		ps:    ps,
		fl:    fl,
	}
}

func (wq *WriteQueue) Run(ctx context.Context) {
	go wq.run(ctx)
}

func (wq *WriteQueue) run(ctx context.Context) {
	for {
		select {
		case qi := <-wq.items:
			if !wq.runItem(qi) {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (wq *WriteQueue) runItem(qi queueItem) bool {
	if err := wq.storeItem(qi); err != nil {
		return wq.fl.OnError(err)
	}

	wq.notifyFlush(qi)
	return true
}

func (wq *WriteQueue) storeItem(qi queueItem) error {
	if qi.sync {
		return wq.ps.Sync(qi.page.Data())
	}

	return wq.ps.Append(qi.lsn, qi.page.Data())
}

func (wq *WriteQueue) notifyFlush(qi queueItem) {
	if qi.lsn.Value() != 0 {
		wq.fl.OnFlush(qi.lsn)
	}
}

func (wq *WriteQueue) Append(lsn record.LogSequenceNumber, p *page.Page) {
	wq.items <- newAppendQueueItem(lsn, p)
}

func (wq *WriteQueue) Sync(lsn record.LogSequenceNumber, p *page.Page) {
	wq.items <- newSyncQueueItem(lsn, p)
}
