package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
)

type syncQueueItem struct {
	page  *page.Page
	reply chan error
}

func newSyncQueueItem(
	page *page.Page,
) *syncQueueItem {
	return &syncQueueItem{
		page:  page,
		reply: make(chan error, 1),
	}
}

func (qi *syncQueueItem) Wait(ctx context.Context) error {
	select {
	case err := <-qi.reply:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (qi *syncQueueItem) Store(ps PageStorage) error {
	err := ps.Sync(qi.page.Data())
	qi.reply <- err
	return err
}

func (qi *syncQueueItem) Page() *page.Page {
	return qi.page
}
