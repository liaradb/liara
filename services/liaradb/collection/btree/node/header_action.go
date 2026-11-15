package node

import (
	"context"

	"github.com/liaradb/liaradb/collection/log"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type headerAction struct {
	frid link.FileRecordID
	data []byte
	prev []byte
}

var _ action = (*headerAction)(nil)

func newHeaderAction(
	frid link.FileRecordID,
	data []byte,
	prev []byte,
) *headerAction {
	return &headerAction{
		frid: frid,
		data: data,
		prev: prev,
	}
}

func (a *headerAction) Log(ctx context.Context, l log.Log) (logpage.LogSequenceNumber, error) {
	return l.UpdateHeader(ctx, a.frid, a.data, a.prev)
}
