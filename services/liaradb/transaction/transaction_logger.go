package transaction

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type transactionLogger struct {
	ctx context.Context
	l   *Logger
}

func newTransactionLogger(
	ctx context.Context,
	l *Logger,
) *transactionLogger {
	return &transactionLogger{
		ctx: ctx,
		l:   l,
	}
}

func (l *transactionLogger) Append(
	slot link.SlotID,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return l.l.Append(l.ctx, slot, data)
}
