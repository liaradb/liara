package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
)

type replaceSyncQueueItem struct {
	page  *page.Page
	reply chan error
}

func newReplaceSyncQueueItem(
	page *page.Page,
) *replaceSyncQueueItem {
	return &replaceSyncQueueItem{
		page:  page,
		reply: make(chan error, 1),
	}
}

func (qi *replaceSyncQueueItem) Wait(ctx context.Context) error {
	select {
	case err := <-qi.reply:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (qi *replaceSyncQueueItem) Store(ps PageStorage) error {
	err := ps.Replace(qi.page.Data())
	qi.reply <- err
	return nil
}

func (qi *replaceSyncQueueItem) Page() *page.Page {
	return qi.page
}
