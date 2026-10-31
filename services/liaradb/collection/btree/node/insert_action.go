package node

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type insertAction struct {
	frid link.FileRecordID
	data []byte
}

var _ action = (*insertAction)(nil)

func newInsertAction(
	frid link.FileRecordID,
	data []byte,
) *insertAction {
	return &insertAction{
		frid: frid,
		data: data,
	}
}

func (a *insertAction) Log(ctx context.Context, l span.Log) (logpage.LogSequenceNumber, error) {
	return l.Insert(ctx, a.frid, a.data)
}
