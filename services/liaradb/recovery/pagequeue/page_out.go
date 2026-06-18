package pagequeue

import (
	"context"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/writequeue"
)

type PageOut struct {
	pool *pagepool.PagePool
	wq   *writequeue.WriteQueue
}

// TODO: This is a duplicate
type PageStorage interface {
	Sync([]byte) error
	Append(record.LogSequenceNumber, []byte) error
	Init([]byte) error
}

func newPageOut(ps PageStorage, pool *pagepool.PagePool) PageOut {
	return PageOut{
		pool: pool,
		wq:   writequeue.New(100, ps, pool),
	}
}

func (po *PageOut) Run(ctx context.Context) error {
	return po.wq.Run(ctx)
}

func (po *PageOut) Append(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	pgs ...*page.Page,
) (*page.Page, bool) {
	// If only one page, do nothing
	l := len(pgs)
	if l <= 1 {
		return nil, false
	}

	for _, p := range pgs[:l-1] {
		po.wq.Append(ctx, lsn, p)
	}

	return pgs[l-1], true
}

func (po *PageOut) Flush(ctx context.Context, current *page.Page) error {
	shadow := po.pool.Get()
	shadow.Fill(current.Data())
	return po.wq.Sync(ctx, shadow)
}
