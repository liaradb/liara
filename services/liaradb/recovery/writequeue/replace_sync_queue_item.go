package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
)

type replaceSyncQueueItem struct {
	page  *logpage.LogPage
	reply chan error
}

func newReplaceSyncQueueItem(
	page *logpage.LogPage,
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

	if err == nil {
		qi.Page().Complete()
	}

	return nil
}

func (qi *replaceSyncQueueItem) Page() *logpage.LogPage {
	return qi.page
}
