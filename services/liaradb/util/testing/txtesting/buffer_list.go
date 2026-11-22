package txtesting

import (
	"context"

	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type bufferList struct {
	storage *storage.Storage
	buffers map[link.FileBlockID]*bufferCount
}

type bufferCount struct {
	buffer *storage.Buffer
	count  int
}

func newBufferList(s *storage.Storage) bufferList {
	return bufferList{
		storage: s,
		buffers: make(map[link.FileBlockID]*bufferCount),
	}
}

func (bl *bufferList) Pin(ctx context.Context, fbid link.FileBlockID) (*storage.Buffer, error) {
	b, err := bl.storage.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (bl *bufferList) PinCurrent(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	b, err := bl.storage.RequestCurrent(ctx, fn)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (bl *bufferList) PinNext(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	b, err := bl.storage.RequestNext(ctx, fn)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (bl *bufferList) Release() {
}

func (bl *bufferList) ReleaseBuffer(b *storage.Buffer) {
	b.Release()
}
