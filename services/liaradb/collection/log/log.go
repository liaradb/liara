package log

import (
	"context"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type Log interface {
	Append(context.Context, link.FileRecordID, []byte) (logpage.LogSequenceNumber, error)
	Insert(context.Context, link.FileRecordID, []byte) (logpage.LogSequenceNumber, error)
	Update(context.Context, link.FileRecordID, []byte, []byte) (logpage.LogSequenceNumber, error)
	Clear(context.Context, link.FileRecordID, []byte) (logpage.LogSequenceNumber, error)
	UpdateHeader(context.Context, link.FileRecordID, []byte, []byte) (logpage.LogSequenceNumber, error)
	// TODO: Use Truncate for splitting B+Tree nodes
	// Truncate(context.Context, link.FileRecordID, []byte) (logpage.LogSequenceNumber, error)
}

type BufferList interface {
	Pin(ctx context.Context, fbid link.FileBlockID) (*storage.Buffer, error)
	PinCurrent(context.Context, link.FileName) (*storage.Buffer, error)
	PinNext(context.Context, link.FileName) (*storage.Buffer, error)
}
