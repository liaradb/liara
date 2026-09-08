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
	c    record.Collection
	log  *log.Log
}

func newLogger(
	tid value.TenantID,
	txid record.TransactionID,
	c record.Collection,
	log *log.Log,
) *Logger {
	return &Logger{
		tid:  tid,
		txid: txid,
		c:    c,
		log:  log,
	}
}

func (l *Logger) Insert(
	ctx context.Context,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Insert(ctx,
		l.tid,
		l.txid,
		time.Now(),
		l.c,
		data)
}

func (l *Logger) Update(
	ctx context.Context,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Update(ctx,
		l.tid,
		l.txid,
		time.Now(),
		l.c,
		data,
		prev)
}
