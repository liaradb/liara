package transaction

import (
	"context"

	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type BufferList struct {
	storage *storage.Storage
	buffers map[link.FileBlockID]*storage.Buffer
}

func NewBufferList(s *storage.Storage) *BufferList {
	return &BufferList{
		storage: s,
		buffers: make(map[link.FileBlockID]*storage.Buffer),
	}
}

func (bl *BufferList) Pin(ctx context.Context, fbid link.FileBlockID) (*storage.Buffer, error) {
	if b, ok := bl.buffers[fbid]; ok {
		return b, nil
	}

	b, err := bl.storage.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	bl.buffers[fbid] = b
	return b, nil
}

func (bl *BufferList) Release() {
	for _, b := range bl.buffers {
		b.Release()
	}
}

func (bl *BufferList) ReleaseBuffer(b *storage.Buffer) {
	b.Release()
	delete(bl.buffers, b.FileBlockID())
}
