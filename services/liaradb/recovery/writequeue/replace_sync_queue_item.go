package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/page"
)

type replaceSyncQueueItem struct {
	page  *page.LogPage
	reply chan error
}

func newReplaceSyncQueueItem(
	page *page.LogPage,
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

func (qi *replaceSyncQueueItem) Page() *page.LogPage {
	return qi.page
}
