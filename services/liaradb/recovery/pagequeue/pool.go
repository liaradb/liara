package pagequeue

import (
	"sync"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type Pool struct {
	pool sync.Pool
}

func NewPool(size int64) Pool {
	return Pool{
		pool: sync.Pool{New: func() any {
			return page.New(size)
		}},
	}
}

func (pl *Pool) Get(pid action.PageID, tlid action.TimeLineID, rl record.Length) *page.Page {
	p := pl.pool.Get().(*page.Page)
	p.Init(pid, tlid, rl)
	return p
}

func (pl *Pool) Put(p *page.Page) {
	pl.pool.Put(p)
}
