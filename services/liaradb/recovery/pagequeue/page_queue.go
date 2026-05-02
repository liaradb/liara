package pagequeue

import (
	"bytes"
	"container/list"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type PageQueue struct {
	pool      page.Pool
	list      list.List
	current   *page.Page
	pid       action.PageID
	tlid      action.TimeLineID
	rl        record.Length
	recordBuf bytes.Buffer
}

func New(size int64) *PageQueue {
	return &PageQueue{
		pool: page.NewPool(size),
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
	data, err := pq.recordToBytes(rc)
	if err != nil {
		return err
	}

	pq.initCurrent()

	if ok := pq.current.Append(data); ok {
		return nil
	}

	pq.next()

	if ok := pq.current.Append(data); !ok {
		return ErrUnableToAppend
	}

	return nil
}

func (pq *PageQueue) initCurrent() {
	if pq.current == nil {
		pq.current = pq.pool.Get(pq.pid, pq.tlid, pq.rl)
	}
}

func (pq *PageQueue) next() {
	pq.list.PushFront(pq.current)
	// TODO: Increment pid
	pq.current = pq.pool.Get(pq.pid, pq.tlid, pq.rl)
}

func (pq *PageQueue) recordToBytes(rc *record.Record) ([]byte, error) {
	// TODO: Use Span instead, as it writes directly
	pq.recordBuf.Reset()
	if err := rc.Write(&pq.recordBuf); err != nil {
		return nil, err
	}

	// We don't need to clone, as the data is copied
	return pq.recordBuf.Bytes(), nil
}

// # Flushing
//   - Flush entire queue to Disk, including Current
func (pq *PageQueue) Flush() error {
	return nil
}
