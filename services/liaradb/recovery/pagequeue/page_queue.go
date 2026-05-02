package pagequeue

import (
	"bytes"
	"container/list"
	"errors"

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
		return errors.New("unable to append")
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

func (wr *PageQueue) recordToBytes(rc *record.Record) ([]byte, error) {
	wr.recordBuf.Reset()
	if err := rc.Write(&wr.recordBuf); err != nil {
		return nil, err
	}

	// We don't need to clone, as the data is copied
	return wr.recordBuf.Bytes(), nil
}
