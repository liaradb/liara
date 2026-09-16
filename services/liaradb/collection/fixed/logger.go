package fixed

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
)

type logger struct {
	ctx        context.Context
	l          *log.Log
	txid       record.TransactionID
	collection record.Collection
}

func NewLogger(
	ctx context.Context,
	l *log.Log,
	txid record.TransactionID,
	collection record.Collection,
) *logger {
	return &logger{
		ctx:        ctx,
		l:          l,
		txid:       txid,
		collection: collection,
	}
}

func (l *logger) Append(rl link.RecordLocator, data []byte) (logpage.LogSequenceNumber, error) {
	return l.l.Insert(l.ctx, l.txid, rl, l.collection, data)
}
