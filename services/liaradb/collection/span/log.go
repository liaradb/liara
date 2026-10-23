package span

import (
	"context"

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
	Append(context.Context, link.RecordLocator, []byte) (logpage.LogSequenceNumber, error)
	Update(context.Context, link.RecordLocator, []byte, []byte) (logpage.LogSequenceNumber, error)
	// Insert(link.RecordLocator, []byte) (logpage.LogSequenceNumber, error) // TODO: Support shifting
	UpdateHeader(context.Context, link.RecordLocator, []byte, []byte) (logpage.LogSequenceNumber, error)
}
