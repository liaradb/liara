package log

import (
	"context"

	"github.com/liaradb/liaradb/log/record"
)

type flushRequest struct {
	lsn      record.LogSequenceNumber
	response chan flushResponse
}

type flushResponse struct {
	err error
}

func newFlushRequest(
	lsn record.LogSequenceNumber,
) *flushRequest {
	return &flushRequest{
		lsn:      lsn,
		response: make(chan flushResponse, 1),
	}
}

func (r *flushRequest) Reply(err error) {
	r.response <- flushResponse{
		err: err,
	}
}

func (r *flushRequest) Wait(ctx context.Context) error {
	select {
	case res := <-r.response:
		return res.err
	case <-ctx.Done():
		return context.Canceled
	}
}
