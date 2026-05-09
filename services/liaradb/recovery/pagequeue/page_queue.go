package pagequeue

import (
	"container/list"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type PageQueue struct {
	pool    Pool
	list    list.List
	current *page.Page
	pid     action.PageID
	tlid    action.TimeLineID
	rl      record.Length
}

func New(size int64) *PageQueue {
	return &PageQueue{
		pool: NewPool(size),
	}
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
	pq.initCurrent()

	t := NewTip(pq.current)
	s := t.Span(int16(rc.Size()))
	if err := rc.Write(s); err != nil {
		return err
	}

	if ok := t.Commit(); !ok {
		return ErrUnableToAppend
	}

	pgs := t.Pages()

	// If pages is empty, do nothing
	l := len(pgs)
	if l == 0 {
		return nil
	}

	pq.list.PushBack(pq.current)
	last := l - 1
	for i, p := range pgs {
		if i == last {
			pq.current = p
		} else {
			pq.list.PushBack(p)
		}
	}

	return nil
}

func (pq *PageQueue) initCurrent() {
	if pq.current == nil {
		pq.current = pq.pool.Get(pq.pid, pq.tlid, pq.rl)
	}
}

// # Flushing
//   - Flush entire queue to Disk, including Current
func (pq *PageQueue) Flush() error {
	// TODO: Implement this
	return nil
}
