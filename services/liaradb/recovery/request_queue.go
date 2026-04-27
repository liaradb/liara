package recovery

import (
	"github.com/liaradb/liaradb/recovery/record"
)

// 1. Append log record
// 2. Fill page
// 3. If page is full, flush
// 4. Request flush to LSN
// 5. Wait for timeout
// 6. Flush to LSN, notify requester
//
// What happens if we flush previous page or segment?
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
