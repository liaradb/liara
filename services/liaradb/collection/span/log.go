package span

import (
	"context"

	"github.com/liaradb/liaradb/collection/tablename"
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
	// TODO: Don't use TableName
	// TODO: Maybe use an interface for string
	XLock(context.Context, tablename.TableName, string) error
	SLock(context.Context, tablename.TableName, string) error

	Append(context.Context, link.RecordID, []byte) (logpage.LogSequenceNumber, error)
	Update(context.Context, link.RecordID, []byte, []byte) (logpage.LogSequenceNumber, error)
	// Insert(link.RecordID, []byte) (logpage.LogSequenceNumber, error) // TODO: Support shifting
	UpdateHeader(context.Context, link.RecordID, []byte, []byte) (logpage.LogSequenceNumber, error)
}
