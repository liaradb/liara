package node

import (
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type headerAction struct {
	rl   link.RecordLocator
	data []byte
	prev []byte
}

var _ action = (*headerAction)(nil)

func newHeaderAction(
	rl link.RecordLocator,
	data []byte,
	prev []byte,
) *headerAction {
	return &headerAction{
		rl:   rl,
		data: data,
		prev: prev,
	}
}

func (a *headerAction) Log(l span.Log) (logpage.LogSequenceNumber, error) {
	return l.UpdateHeader(a.rl, a.data, a.prev)
}
