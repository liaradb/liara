package log

import (
	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/log/record"
)

type flushRequest = async.Request[record.LogSequenceNumber, struct{}]

func newFlushRequest(
	lsn record.LogSequenceNumber,
) *flushRequest {
	return async.NewRequest[record.LogSequenceNumber, struct{}](lsn)
}
