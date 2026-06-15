package pagequeue

import (
	"container/list"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/iterator"
)

type PageOut struct {
	pool   *Pool
	shadow *page.Page
	list   list.List
	ps     PageStorage
	lsn    record.LogSequenceNumber
}

type PageStorage interface {
	Sync([]byte) error
	Append(record.LogSequenceNumber, []byte) error
	Init([]byte) error
}

func newPageOut(ps PageStorage, pool *Pool) PageOut {
	return PageOut{
		pool:   pool,
		ps:     ps,
		shadow: pool.Get(),
	}
}

func (po *PageOut) Count() int {
	return po.list.Len() + 1
}

func (po *PageOut) Append(
	lsn record.LogSequenceNumber,
	c *page.Page,
	pgs []*page.Page,
) (*page.Page, bool) {
	po.setLSN(lsn)

	// If pages is empty, do nothing
	l := len(pgs)
	if l == 0 {
		return nil, false
	}

	po.push(c)

	last := l - 1
	for i, p := range pgs {
		if i == last {
			// TODO: Should we do something with p?
			return p, true
		} else {
			po.push(p)
		}
	}

	// TODO: This should never happen
	return nil, false
}

func (po *PageOut) setLSN(lsn record.LogSequenceNumber) {
	po.lsn = lsn
}

func (po *PageOut) push(p *page.Page) {
	po.list.PushBack(p)
}

// # Flushing
//   - Flush entire queue to Disk, including Current
func (po *PageOut) Flush(current *page.Page) error {
	po.syncShadow(current)

	if po.list.Len() == 0 {
		return po.syncCurrent()
	}

	// TODO: Implement error handling
	i := 0
	for p := range iterator.Forward[*page.Page](&po.list) {
		if i == 0 { // Sync first
			if err := po.syncPage(p); err != nil {
				return err
			}
		} else {
			if err := po.appendPage(p); err != nil {
				return err
			}
		}

		po.pool.Put(p)
		i++
	}

	po.list.Init()

	return po.appendCurrent()
}

func (po *PageOut) syncShadow(current *page.Page) {
	po.shadow.Fill(current.Data())
}

func (po *PageOut) syncCurrent() error {
	return po.syncPage(po.shadow)
}

func (po *PageOut) syncPage(p *page.Page) error {
	return po.ps.Sync(p.Data())
}

func (po *PageOut) appendCurrent() error {
	return po.appendPage(po.shadow)
}

func (po *PageOut) appendPage(p *page.Page) error {
	return po.ps.Append(po.lsn, p.Data())
}

func (po *PageOut) Clear() {
	for p := range iterator.Forward[*page.Page](&po.list) {
		po.pool.Put(p)
	}
	po.list.Init()
}
