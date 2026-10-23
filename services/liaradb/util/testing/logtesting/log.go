package logtesting

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type MockLog struct {
}

func (t *MockLog) Append(context.Context, link.RecordLocator, []byte) (logpage.LogSequenceNumber, error) {
	return logpage.LogSequenceNumber{}, nil
}

func (t *MockLog) Update(context.Context, link.RecordLocator, []byte, []byte) (logpage.LogSequenceNumber, error) {
	return logpage.LogSequenceNumber{}, nil
}

func (t *MockLog) UpdateHeader(context.Context, link.RecordLocator, []byte, []byte) (logpage.LogSequenceNumber, error) {
	return logpage.LogSequenceNumber{}, nil
}
