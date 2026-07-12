package recordqueue

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/record"
)

type AppendRequest = async.Request[AppendValue, record.LogSequenceNumber]

type AppendHandler struct {
	reqs async.Handler[AppendValue, record.LogSequenceNumber]
}

func NewAppendHandler() AppendHandler {
	return AppendHandler{
		reqs: make(async.Handler[AppendValue, record.LogSequenceNumber]),
	}
}

func (h AppendHandler) Reqs() async.Handler[AppendValue, record.LogSequenceNumber] {
	return h.reqs
}

func (h AppendHandler) Append(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	time time.Time,
	action record.Action,
	collection record.Collection,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	return h.reqs.Send(ctx, AppendValue{
		tid:        tid,
		txid:       txid,
		time:       time,
		action:     action,
		collection: collection,
		data:       data,
		reverse:    reverse,
	})
}

func (h AppendHandler) AppendAndWait(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	time time.Time,
	action record.Action,
	collection record.Collection,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	return h.reqs.Send(ctx, AppendValue{
		tid:        tid,
		txid:       txid,
		time:       time,
		action:     action,
		collection: collection,
		data:       data,
		reverse:    reverse,
		wait:       true,
	})
}

type AppendValue struct {
	tid        value.TenantID
	txid       record.TransactionID
	time       time.Time
	action     record.Action
	collection record.Collection
	data       []byte
	reverse    []byte
	wait       bool
}

func (av *AppendValue) Record(lsn record.LogSequenceNumber) *record.Record {
	return record.New(
		lsn,
		av.tid,
		av.txid,
		record.NewTime(av.time),
		av.action,
		av.collection,
		av.data, av.reverse)
}

func (av AppendValue) IsWait() bool {
	return av.wait
}
