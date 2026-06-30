package pagepool

import (
	"sync"

	"github.com/liaradb/liaradb/recovery/logpage"
)

// TODO: Can we implement this without sync.PagePool?
type PagePool struct {
	pool sync.Pool
}

func New(size int16, slotHeaderSize int16) PagePool {
	return PagePool{
		pool: sync.Pool{New: func() any {
			return logpage.New(size, slotHeaderSize)
		}},
	}
}

func (pl *PagePool) Get() *logpage.LogPage {
	p := pl.pool.Get().(*logpage.LogPage)
	p.Clear()
	return p
}

func (pl *PagePool) Put(p *logpage.LogPage) {
	p.Reset()
	pl.pool.Put(p)
}
