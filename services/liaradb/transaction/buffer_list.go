package transaction

import (
	"context"

	"github.com/liaradb/liaradb/storage"
)

type BufferList struct {
	storage *storage.Storage
	buffers map[storage.BlockID]*storage.BufferPage
}

func NewBufferList(s *storage.Storage) *BufferList {
	return &BufferList{
		storage: s,
		buffers: make(map[storage.BlockID]*storage.BufferPage),
	}
}

func (bl *BufferList) Pin(ctx context.Context, bid storage.BlockID) (*storage.BufferPage, error) {
	if b, ok := bl.buffers[bid]; ok {
		return b, nil
	}

	b, err := bl.storage.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	bp := storage.NewBufferPage(b)

	bl.buffers[bid] = bp
	return bp, nil
}

func (bl *BufferList) Release() {
	for _, b := range bl.buffers {
		b.Release()
	}
}

func (bl *BufferList) ReleaseBuffer(b *storage.BufferPage) {
	b.Release()
	delete(bl.buffers, b.BlockID())
}
