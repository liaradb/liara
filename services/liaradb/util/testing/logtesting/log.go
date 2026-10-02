package logtesting

import (
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type MockLog struct {
}

func (t *MockLog) Append(link.RecordLocator, []byte) (logpage.LogSequenceNumber, error) {
	return logpage.LogSequenceNumber{}, nil
}

func (t *MockLog) UpdateHeader(link.RecordLocator, []byte, []byte) (logpage.LogSequenceNumber, error) {
	return logpage.LogSequenceNumber{}, nil
}
