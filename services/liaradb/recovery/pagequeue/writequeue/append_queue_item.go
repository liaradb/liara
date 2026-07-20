package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
)

type appendQueueItem struct {
	lsn  logpage.LogSequenceNumber
	page *logpage.LogPage
}

func newAppendQueueItem(
	lsn logpage.LogSequenceNumber,
	page *logpage.LogPage,
) *appendQueueItem {
	return &appendQueueItem{
		lsn:  lsn,
		page: page,
	}
}

func (*appendQueueItem) Wait(context.Context) error {
	return nil
}

func (qi *appendQueueItem) Store(ps PageStorage) error {
	if err := ps.Append(qi.lsn, qi.page.Data()); err != nil {
		return err
	}

	qi.Page().Complete()
	return nil
}

func (qi *appendQueueItem) Page() *logpage.LogPage {
	return qi.page
}
