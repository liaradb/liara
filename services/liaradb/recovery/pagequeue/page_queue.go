package pagequeue

import (
	"container/list"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/iterator"
)

type PageQueue struct {
	pool    Pool
	list    list.List
	current *page.Page
	shadow  *page.Page
	ps      PageStorage
	lsn     record.LogSequenceNumber
}

type PageStorage interface {
	Sync([]byte) error
	Append(record.LogSequenceNumber, []byte) error
	Init([]byte) error
}

func New(
	ps PageStorage,
	size int16,
	headerSize int16,
	slotHeaderSize int16,
) *PageQueue {
	pq := &PageQueue{
		pool: NewPool(size, headerSize, slotHeaderSize),
		ps:   ps,
	}

	pq.init()

	return pq
}

func (pq *PageQueue) init() {
	pq.shadow = pq.pool.Get()
	pq.current = pq.pool.Get()
}

func (pq *PageQueue) Init(data []byte) {
	pq.current.Fill(data)
}

func (pq *PageQueue) Count() int {
	return pq.list.Len() + 1
}

// # Append
//   - Compare size of new Record to remaining space in current Page
//   - If it fits, append to the current Page
//   - If it spans, generate a new list of Pages to fit
//   - Append Record as Span to the list
//   - Append list to queue, up to but not including, current
//   - If current Page is entirely full, append current to list and swap current for next Page
func (pq *PageQueue) Append(rc *record.Record) error {
	t := NewTip(&pq.pool, pq.current)
	s := t.Span(int16(rc.Size()))
	if err := rc.Write(s); err != nil {
		return err
	}

	pgs, ok := t.Commit()
	if !ok {
		return ErrUnableToAppend
	}

	pq.lsn = rc.LogSequenceNumber()
	pq.appendPages(pgs)
	return nil
}

func (pq *PageQueue) appendPages(pgs []*page.Page) {
	// If pages is empty, do nothing
	l := len(pgs)
	if l == 0 {
		return
	}

	pq.list.PushBack(pq.current)
	last := l - 1
	for i, p := range pgs {
		if i == last {
			pq.replaceCurrent(p)
		} else {
			pq.list.PushBack(p)
		}
	}
}

// # Flushing
//   - Flush entire queue to Disk, including Current
func (pq *PageQueue) Flush() error {
	pq.syncShadow()

	if pq.list.Len() == 0 {
		return pq.syncCurrent()
	}

	// TODO: Implement error handling
	i := 0
	for p := range iterator.Forward[*page.Page](&pq.list) {
		if i == 0 { // Sync first
			if err := pq.syncPage(p); err != nil {
				return err
			}
		} else {
			if err := pq.appendPage(p); err != nil {
				return err
			}
		}

		pq.pool.Put(p)
		i++
	}

	pq.list.Init()

	return pq.appendCurrent()
}

func (pq *PageQueue) replaceCurrent(p *page.Page) {
	pq.current = p
}

func (pq *PageQueue) syncShadow() {
	pq.shadow.Fill(pq.current.Data())
}

func (pq *PageQueue) syncCurrent() error {
	return pq.syncPage(pq.shadow)
}

func (pq *PageQueue) syncPage(p *page.Page) error {
	return pq.ps.Sync(p.Data())
}

func (pq *PageQueue) appendCurrent() error {
	return pq.appendPage(pq.shadow)
}

func (pq *PageQueue) appendPage(p *page.Page) error {
	return pq.ps.Append(pq.lsn, p.Data())
}

func (pq *PageQueue) Clear() {
	for p := range iterator.Forward[*page.Page](&pq.list) {
		pq.pool.Put(p)
	}
	pq.list.Init()
}
