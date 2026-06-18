package pagepool

import (
	"sync"

	"github.com/liaradb/liaradb/encoder/page"
)

// TODO: Can we implement this without sync.PagePool?
type PagePool struct {
	pool sync.Pool
}

func New(size int16, headerSize int16, slotHeaderSize int16) PagePool {
	return PagePool{
		pool: sync.Pool{New: func() any {
			return page.New(size, headerSize, slotHeaderSize)
		}},
	}
}

func (pl *PagePool) Get() *page.Page {
	p := pl.pool.Get().(*page.Page)
	p.Clear()
	return p
}

func (pl *PagePool) Put(p *page.Page) {
	pl.pool.Put(p)
}
