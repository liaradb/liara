package recordqueue

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/util/async"
)

type AppendRequest[R pagequeue.Record] = async.Request[AppendValue[R], logpage.LogSequenceNumber]

type AppendHandler[R pagequeue.Record] struct {
	reqs async.Handler[AppendValue[R], logpage.LogSequenceNumber]
}

func NewAppendHandler[R pagequeue.Record]() AppendHandler[R] {
	return AppendHandler[R]{
		reqs: make(async.Handler[AppendValue[R], logpage.LogSequenceNumber]),
	}
}

func (h AppendHandler[R]) Reqs() async.Handler[AppendValue[R], logpage.LogSequenceNumber] {
	return h.reqs
}

func (h AppendHandler[R]) Append(
	ctx context.Context,
	record R,
) (logpage.LogSequenceNumber, error) {
	return h.reqs.Send(ctx, AppendValue[R]{
		record: record,
	})
}

func (h AppendHandler[R]) AppendAndWait(
	ctx context.Context,
	record R,
) (logpage.LogSequenceNumber, error) {
	return h.reqs.Send(ctx, AppendValue[R]{
		record: record,
		wait:   true,
	})
}

type AppendValue[R pagequeue.Record] struct {
	record R
	wait   bool
}

func (av *AppendValue[R]) Record(lsn logpage.LogSequenceNumber) R {
	return av.record
}

func (av AppendValue[R]) IsWait() bool {
	return av.wait
}
