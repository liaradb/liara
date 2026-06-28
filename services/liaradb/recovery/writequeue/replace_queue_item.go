package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
)

type replaceQueueItem struct {
	page *logpage.LogPage
}

func newReplaceQueueItem(
	page *logpage.LogPage,
) *replaceQueueItem {
	return &replaceQueueItem{
		page: page,
	}
}

func (*replaceQueueItem) Wait(context.Context) error {
	return nil
}

func (qi *replaceQueueItem) Store(ps PageStorage, fl Flusher) error {
	if err := ps.Replace(qi.page.Data()); err != nil {
		return err
	}

	fl.OnFlush(qi.Page().LSN())
	return nil
}

func (qi *replaceQueueItem) Page() *logpage.LogPage {
	return qi.page
}
