package logger

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
)

type Logger struct {
	log  *log.Log
	txid record.TransactionID
}

var _ span.Log = (*Logger)(nil)

func New(
	log *log.Log,
	txid record.TransactionID,
) *Logger {
	return &Logger{
		log:  log,
		txid: txid,
	}
}

func (l *Logger) Append(
	ctx context.Context,
	rid link.RecordID,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Insert(ctx,
		l.txid,
		rid,
		data)
}

func (l *Logger) Update(
	ctx context.Context,
	rid link.RecordID,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Update(ctx,
		l.txid,
		rid,
		data,
		prev)
}

func (l *Logger) UpdateHeader(
	ctx context.Context,
	rid link.RecordID,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.UpdateHeader(ctx,
		l.txid,
		rid,
		data,
		prev)
}
