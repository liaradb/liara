package node

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type insertAction struct {
	rid  link.RecordID
	data []byte
}

var _ action = (*insertAction)(nil)

func newInsertAction(
	rid link.RecordID,
	data []byte,
) *insertAction {
	return &insertAction{
		rid:  rid,
		data: data,
	}
}

func (a *insertAction) Log(ctx context.Context, l span.Log) (logpage.LogSequenceNumber, error) {
	return l.Insert(ctx, a.rid, a.data)
}
