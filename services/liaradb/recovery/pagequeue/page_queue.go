package pagequeue

import (
	"container/list"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/iterator"
)

type PageQueue struct {
	size           int16
	headerSize     int16
	slotHeaderSize int16
	pool           Pool
	list           list.List
	current        *page.Page
	pid            action.PageID
	tlid           action.TimeLineID
	rl             record.Length
}

func New(size int16, headerSize int16, slotHeaderSize int16) *PageQueue {
	pq := &PageQueue{
		size:           size,
		headerSize:     headerSize,
		slotHeaderSize: slotHeaderSize,
		pool:           NewPool(size, headerSize, slotHeaderSize),
	}

	pq.initCurrent()

	return pq
}

func (pq *PageQueue) initCurrent() {
	if pq.current == nil {
		pq.current = pq.pool.Get()
	}
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
			pq.current = p
		} else {
			pq.list.PushBack(p)
		}
	}
}

// # Flushing
//   - Flush entire queue to Disk, including Current
func (pq *PageQueue) Flush() error {
	// TODO: Implement this
	for p := range iterator.Forward[*page.Page](&pq.list) {
		// Flush p
		pq.pool.Put(p)
	}
	pq.list.Init()
	// Flush current
	return nil
}

func (pq *PageQueue) Pages() iter.Seq[*page.Page] {
	return func(yield func(*page.Page) bool) {
		for p := range iterator.Forward[*page.Page](&pq.list) {
			if !yield(p) {
				return
			}
		}
		yield(pq.current)
	}
}

func (pq *PageQueue) Clear() {
	for p := range iterator.Forward[*page.Page](&pq.list) {
		pq.pool.Put(p)
	}
	pq.list.Init()
}
