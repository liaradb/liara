package node

import (
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type AppendAction struct {
	rl   link.RecordLocator
	data []byte
}

var _ Action = (*AppendAction)(nil)

func (a *AppendAction) Log(l span.Log) (logpage.LogSequenceNumber, error) {
	return l.Append(a.rl, a.data)
}
