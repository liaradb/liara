package log

import (
	"time"

	"github.com/liaradb/liaradb/log/record"
)

type appendRequest struct {
	tid      record.TransactionID
	time     time.Time
	action   record.Action
	data     []byte
	reverse  []byte
	response chan appendResponse
}

type appendResponse struct {
	lsn record.LogSequenceNumber
	err error
}

func newAppendRequest(
	tid record.TransactionID,
	time time.Time,
	action record.Action,
	data []byte,
	reverse []byte,
) *appendRequest {
	return &appendRequest{
		tid:      tid,
		time:     time,
		action:   action,
		data:     data,
		reverse:  reverse,
		response: make(chan appendResponse, 1),
	}
}

func (r *appendRequest) Reply(lsn record.LogSequenceNumber, err error) {
	r.response <- appendResponse{
		lsn: lsn,
		err: err,
	}
}
