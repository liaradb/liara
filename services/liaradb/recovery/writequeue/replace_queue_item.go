package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/page"
)

type replaceQueueItem struct {
	page *page.LogPage
}

func newReplaceQueueItem(
	page *page.LogPage,
) *replaceQueueItem {
	return &replaceQueueItem{
		page: page,
	}
}

func (*replaceQueueItem) Wait(context.Context) error {
	return nil
}

func (qi *replaceQueueItem) Store(ps PageStorage) error {
	return ps.Replace(qi.page.Data())
}

func (qi *replaceQueueItem) Page() *page.LogPage {
	return qi.page
}
