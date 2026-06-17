package pagequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type appendQueueItem struct {
	lsn  record.LogSequenceNumber
	page *page.Page
}

func newAppendQueueItem(
	lsn record.LogSequenceNumber,
	page *page.Page,
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
