package logger

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
)

// TODO: Simplify transaction package reference
type Logger struct {
	ctx  context.Context
	log  *log.Log
	txid record.TransactionID
	c    record.Collection
}

var _ span.Log = (*Logger)(nil)

func New(
	ctx context.Context,
	log *log.Log,
	txid record.TransactionID,
	c record.Collection,
) *Logger {
	return &Logger{
		ctx:  ctx,
		log:  log,
		txid: txid,
		c:    c,
	}
}

func (l *Logger) Append(
	rl link.RecordLocator,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Insert(l.ctx,
		l.txid,
		rl,
		l.c,
		data)
}

func (l *Logger) Update(
	rl link.RecordLocator,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Update(l.ctx,
		l.txid,
		rl,
		l.c,
		data,
		prev)
}
