package fixed

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
)

type logger struct {
	ctx context.Context
	l   *log.Log
}

func NewLogger(
	ctx context.Context,
	l *log.Log,
) *logger {
	return &logger{
		ctx: ctx,
		l:   l,
	}
}

func (l *logger) Append(link.SlotID, []byte) (logpage.LogSequenceNumber, error) {
	return logpage.LogSequenceNumber{}, nil
}
