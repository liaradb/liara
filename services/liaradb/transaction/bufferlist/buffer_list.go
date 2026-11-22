package bufferlist

import (
	"context"

	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type BufferList struct {
	storage *storage.Storage
	buffers map[link.FileBlockID]*bufferCount
}

type bufferCount struct {
	buffer *storage.Buffer
	count  int
}

func NewBufferList(s *storage.Storage) *BufferList {
	return &BufferList{
		storage: s,
		buffers: make(map[link.FileBlockID]*bufferCount),
	}
}

func (bl *BufferList) Pin(ctx context.Context, fbid link.FileBlockID) (*storage.Buffer, error) {
	b, err := bl.storage.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (bl *BufferList) PinCurrent(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	b, err := bl.storage.RequestCurrent(ctx, fn)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (bl *BufferList) PinNext(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	b, err := bl.storage.RequestNext(ctx, fn)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (bl *BufferList) Release() {
}

func (bl *BufferList) ReleaseBuffer(b *storage.Buffer) {
	b.Release()
}
