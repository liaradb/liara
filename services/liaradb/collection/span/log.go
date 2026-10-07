package span

import (
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

// # Operations
//   - Set
//   - Insert
//   - Push
//   - Pop
//   - Set Range
//   - Delete Range
//   - Change header
type Log interface {
	Append(link.RecordLocator, []byte) (logpage.LogSequenceNumber, error)
	// Insert(link.RecordLocator, []byte) (logpage.LogSequenceNumber, error) // TODO: Support shifting
	UpdateHeader(link.RecordLocator, []byte, []byte) (logpage.LogSequenceNumber, error)
}
