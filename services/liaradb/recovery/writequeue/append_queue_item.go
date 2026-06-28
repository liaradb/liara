package writequeue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/record"
)

type appendQueueItem struct {
	lsn  record.LogSequenceNumber
	page *logpage.LogPage
}

func newAppendQueueItem(
	lsn record.LogSequenceNumber,
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
	return ps.Append(qi.lsn, qi.page.Data())
}

func (qi *appendQueueItem) Page() *logpage.LogPage {
	return qi.page
}
