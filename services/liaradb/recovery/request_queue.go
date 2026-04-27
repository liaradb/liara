package recovery

import (
	"github.com/liaradb/liaradb/recovery/record"
)

type requestQueue struct {
	items []*flushRequest
}

func (rq *requestQueue) Add(r *flushRequest) {
	rq.items = append(rq.items, r)
}

func (rq *requestQueue) SendUpToLSN(lsn record.LogSequenceNumber) {
	for _, r := range rq.items {
		r.Reply(nil)
	}

	clear(rq.items)
	rq.items = rq.items[:0]
}
