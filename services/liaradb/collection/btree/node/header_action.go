package node

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type headerAction struct {
	rid  link.RecordID
	data []byte
	prev []byte
}

var _ action = (*headerAction)(nil)

func newHeaderAction(
	rid link.RecordID,
	data []byte,
	prev []byte,
) *headerAction {
	return &headerAction{
		rid:  rid,
		data: data,
		prev: prev,
	}
}

func (a *headerAction) Log(ctx context.Context, l span.Log) (logpage.LogSequenceNumber, error) {
	return l.UpdateHeader(ctx, a.rid, a.data, a.prev)
}
