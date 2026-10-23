package node

import (
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type InsertAction struct {
	rl   link.RecordLocator
	data []byte
}

var _ Action = (*InsertAction)(nil)

func (a *InsertAction) Log(l span.Log) (logpage.LogSequenceNumber, error) {
	return l.Append(a.rl, a.data)
}
