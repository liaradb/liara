package transaction

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
)

type Logger struct {
	tid  value.TenantID
	txid record.TransactionID
	l    *log.Log
}

func (l *Logger) Insert(
	ctx context.Context,
	now time.Time,
	collection record.Collection,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return l.l.Insert(ctx, l.tid, l.txid, now, collection, data)
}

func (l *Logger) Update(
	ctx context.Context,
	now time.Time,
	collection record.Collection,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return l.l.Update(ctx, l.tid, l.txid, now, collection, data, prev)
}
