package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
)

type replaceQueueItem struct {
	page *page.Page
}

func newReplaceQueueItem(
	page *page.Page,
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

func (qi *replaceQueueItem) Page() *page.Page {
	return qi.page
}
