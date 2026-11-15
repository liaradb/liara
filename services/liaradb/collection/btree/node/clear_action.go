package node

import (
	"context"

	"github.com/liaradb/liaradb/collection/log"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type clearAction struct {
	frid link.FileRecordID
	prev []byte
}

var _ action = (*clearAction)(nil)

func newClearAction(
	frid link.FileRecordID,
	prev []byte,
) *clearAction {
	return &clearAction{
		frid: frid,
		prev: prev,
	}
}

func (a *clearAction) Log(ctx context.Context, l log.Log) (logpage.LogSequenceNumber, error) {
	return l.Clear(ctx, a.frid, a.prev)
}
