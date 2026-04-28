package segment

import (
	"sync"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type PagePool struct {
	pool sync.Pool
}

func NewPagePool(size int64) *PagePool {
	return &PagePool{
		pool: sync.Pool{New: func() any {
			return page.New(size)
		}},
	}
}

func (pp *PagePool) Get(pid action.PageID, tlid action.TimeLineID, rl record.Length) *page.Page {
	p := pp.pool.Get().(*page.Page)
	p.Init(pid, tlid, rl)
	return p
}

func (pp *PagePool) Return(p *page.Page) {
	pp.pool.Put(p)
}
