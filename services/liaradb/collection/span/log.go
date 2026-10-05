package span

import (
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type Log interface {
	Append(link.RecordLocator, []byte) (logpage.LogSequenceNumber, error)
	UpdateHeader(link.RecordLocator, []byte, []byte) (logpage.LogSequenceNumber, error)
}
