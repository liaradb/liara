package node

import (
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type insertAction struct {
	rl   link.RecordLocator
	data []byte
}

var _ action = (*insertAction)(nil)

func newInsertAction(
	rl link.RecordLocator,
	data []byte,
) *insertAction {
	return &insertAction{
		rl:   rl,
		data: data,
	}
}

func (a *insertAction) Log(l span.Log) (logpage.LogSequenceNumber, error) {
	return l.Append(a.rl, a.data)
}
