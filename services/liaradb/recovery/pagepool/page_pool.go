package pagepool

import (
	"sync"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/record"
)

// TODO: Can we implement this without sync.PagePool?
type PagePool struct {
	pool sync.Pool
}

func New(size int16, slotHeaderSize int16) PagePool {
	return PagePool{
		pool: sync.Pool{New: func() any {
			p := page.New(size, logpage.HeaderSize, slotHeaderSize)
			return &logpage.LogPage{
				Page: p,
			}
		}},
	}
}

func (pl *PagePool) Get() *logpage.LogPage {
	p := pl.pool.Get().(*logpage.LogPage)
	p.Page.Clear()
	return p
}

func (pl *PagePool) GetLsn(lsn record.LogSequenceNumber) *logpage.LogPage {
	p := pl.pool.Get().(*logpage.LogPage)
	p.Clear(lsn)
	return p
}

func (pl *PagePool) Put(p *logpage.LogPage) {
	p.Reset()
	pl.pool.Put(p)
}
