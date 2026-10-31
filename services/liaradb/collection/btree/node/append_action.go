package node

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type appendAction struct {
	frid link.FileRecordID
	data []byte
}

var _ action = (*appendAction)(nil)

func newAppendAction(
	frid link.FileRecordID,
	data []byte,
) *appendAction {
	return &appendAction{
		frid: frid,
		data: data,
	}
}

func (a *appendAction) Log(ctx context.Context, l span.Log) (logpage.LogSequenceNumber, error) {
	return l.Append(ctx, a.frid, a.data)
}
