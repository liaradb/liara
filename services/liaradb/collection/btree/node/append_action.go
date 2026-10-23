package node

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type appendAction struct {
	rl   link.RecordLocator
	data []byte
}

var _ action = (*appendAction)(nil)

func newAppendAction(
	rl link.RecordLocator,
	data []byte,
) *appendAction {
	return &appendAction{
		rl:   rl,
		data: data,
	}
}

func (a *appendAction) Log(ctx context.Context, l span.Log) (logpage.LogSequenceNumber, error) {
	return l.Append(ctx, a.rl, a.data)
}
