package recovery

import (
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/record"
)

type appendRequest = async.Request[appendValue, record.LogSequenceNumber]

type appendValue struct {
	tid        value.TenantID
	txid       record.TransactionID
	time       time.Time
	action     record.Action
	collection record.Collection
	data       []byte
	reverse    []byte
}

func (av *appendValue) Record(lsn record.LogSequenceNumber) *record.Record {
	return record.New(
		lsn,
		av.tid,
		av.txid,
		record.NewTime(av.time),
		av.action,
		av.collection,
		av.data, av.reverse)
}
