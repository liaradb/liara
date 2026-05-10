package pagequeue

import (
	"sync"

	"github.com/liaradb/liaradb/encoder/page"
)

type Pool struct {
	pool sync.Pool
}

func NewPool(size int16, headerSize int16, slotHeaderSize int16) Pool {
	return Pool{
		pool: sync.Pool{New: func() any {
			return page.New(size, headerSize, slotHeaderSize)
		}},
	}
}

func (pl *Pool) Get() *page.Page {
	p := pl.pool.Get().(*page.Page)
	p.Reset()
	return p
}

func (pl *Pool) Put(p *page.Page) {
	pl.pool.Put(p)
}
