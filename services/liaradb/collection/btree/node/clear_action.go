package node

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type clearAction struct {
	rid  link.RecordID
	prev []byte
}

var _ action = (*clearAction)(nil)

func newClearAction(
	rl link.RecordID,
	prev []byte,
) *clearAction {
	return &clearAction{
		rid:  rl,
		prev: prev,
	}
}

func (a *clearAction) Log(ctx context.Context, l span.Log) (logpage.LogSequenceNumber, error) {
	return l.Clear(ctx, a.rid, a.prev)
}
