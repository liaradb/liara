package pagequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type queueItem struct {
	lsn   record.LogSequenceNumber
	page  *page.Page
	sync  bool
	reply chan struct{}
}

func newSyncQueueItem(
	lsn record.LogSequenceNumber,
	page *page.Page,
) queueItem {
	return queueItem{
		lsn:   lsn,
		page:  page,
		sync:  true,
		reply: make(chan struct{}, 1),
	}
}

func newAppendQueueItem(
	lsn record.LogSequenceNumber,
	page *page.Page,
) queueItem {
	return queueItem{
		lsn:  lsn,
		page: page,
	}
}

func (qi *queueItem) Reply() {
	if qi.sync {
		qi.reply <- struct{}{}
	}
}

func (qi *queueItem) Wait(ctx context.Context) {
	if qi.sync {
		qi.wait(ctx)
	}
}

func (qi *queueItem) wait(ctx context.Context) {
	select {
	case <-qi.reply:
	case <-ctx.Done():
	}
}
