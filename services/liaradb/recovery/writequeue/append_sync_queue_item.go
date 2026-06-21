package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type appendSyncQueueItem struct {
	lsn   record.LogSequenceNumber
	page  *page.Page
	reply chan error
}

func newAppendSyncQueueItem(
	lsn record.LogSequenceNumber,
	page *page.Page,
) *appendSyncQueueItem {
	return &appendSyncQueueItem{
		lsn:   lsn,
		page:  page,
		reply: make(chan error, 1),
	}
}

func (qi *appendSyncQueueItem) Wait(ctx context.Context) error {
	select {
	case err := <-qi.reply:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (qi *appendSyncQueueItem) Store(ps PageStorage) error {
	err := ps.Append(qi.lsn, qi.page.Data())
	qi.reply <- err
	return nil
}

func (qi *appendSyncQueueItem) Page() *page.Page {
	return qi.page
}
