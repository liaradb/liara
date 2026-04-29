package page

import (
	"sync"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

type Pool struct {
	pool sync.Pool
}

func NewPool(size int64) *Pool {
	return &Pool{
		pool: sync.Pool{New: func() any {
			return New(size)
		}},
	}
}

func (pl *Pool) Get(pid action.PageID, tlid action.TimeLineID, rl record.Length) *Page {
	p := pl.pool.Get().(*Page)
	p.Init(pid, tlid, rl)
	return p
}

func (pl *Pool) Return(p *Page) {
	pl.pool.Put(p)
}
