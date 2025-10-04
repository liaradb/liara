package log

import (
	"time"

	"github.com/liaradb/liaradb/log/record"
)

type request struct {
	tid      record.TransactionID
	time     time.Time
	action   record.Action
	data     []byte
	reverse  []byte
	response chan response
}

type response struct {
	lsn record.LogSequenceNumber
	err error
}

func newRequest(
	tid record.TransactionID,
	time time.Time,
	action record.Action,
	data []byte,
	reverse []byte,
) *request {
	return &request{
		tid:      tid,
		time:     time,
		action:   action,
		data:     data,
		reverse:  reverse,
		response: make(chan response, 1),
	}
}

func (r *request) Reply(lsn record.LogSequenceNumber, err error) {
	r.response <- response{
		lsn: lsn,
		err: err,
	}
}
