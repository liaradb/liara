package log

import (
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/log/record"
)

type appendRequest = async.Request[appendValue, record.LogSequenceNumber]

type appendValue struct {
	tid     record.TransactionID
	time    time.Time
	action  record.Action
	data    []byte
	reverse []byte
}

func newAppendRequest(
	tid record.TransactionID,
	time time.Time,
	action record.Action,
	data []byte,
	reverse []byte,
) *appendRequest {
	return async.NewRequest[appendValue, record.LogSequenceNumber](appendValue{
		tid:     tid,
		time:    time,
		action:  action,
		data:    data,
		reverse: reverse,
	})
}
