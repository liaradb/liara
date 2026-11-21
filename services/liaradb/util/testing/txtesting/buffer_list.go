package txtesting

import (
	"context"

	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type bufferList struct {
	storage *storage.Storage
	buffers map[link.FileBlockID]*storage.Buffer
}

func newBufferList(s *storage.Storage) bufferList {
	return bufferList{
		storage: s,
		buffers: make(map[link.FileBlockID]*storage.Buffer),
	}
}

func (bl *bufferList) Pin(ctx context.Context, fbid link.FileBlockID) (*storage.Buffer, error) {
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

func (bl *bufferList) PinCurrent(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	b, err := bl.storage.RequestCurrent(ctx, fn)
	if err != nil {
		return nil, err
	}

	bl.buffers[b.FileBlockID()] = b
	return b, nil
}

func (bl *bufferList) PinNext(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	b, err := bl.storage.RequestNext(ctx, fn)
	if err != nil {
		return nil, err
	}

	bl.buffers[b.FileBlockID()] = b
	return b, nil
}

func (bl *bufferList) Release() {
	for _, b := range bl.buffers {
		b.Release()
	}
}

func (bl *bufferList) ReleaseBuffer(b *storage.Buffer) {
	if !b.Dirty() {
		b.Release()
		delete(bl.buffers, b.FileBlockID())
	}
}
