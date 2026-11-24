package bufferlist

import (
	"context"

	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type BufferList struct {
	storage *storage.Storage
	buffers map[link.FileBlockID]*Buffer
}

func NewBufferList(s *storage.Storage) *BufferList {
	return &BufferList{
		storage: s,
		buffers: make(map[link.FileBlockID]*Buffer),
	}
}

func (bl *BufferList) Pin(ctx context.Context, fbid link.FileBlockID) (*Buffer, error) {
	if b, ok := bl.buffers[fbid]; ok {
		return b, nil
	}

	b, err := bl.storage.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	bf := NewBuffer(bl, b)
	bl.buffers[fbid] = bf

	return bf, nil
}

func (bl *BufferList) PinCurrent(ctx context.Context, fn link.FileName) (*Buffer, error) {
	b, err := bl.storage.RequestCurrent(ctx, fn)
	if err != nil {
		return nil, err
	}

	if bf, ok := bl.buffers[b.FileBlockID()]; ok {
		b.Release()
		return bf, nil
	}

	bf := NewBuffer(bl, b)
	bl.buffers[b.FileBlockID()] = bf
	return bf, nil
}

func (bl *BufferList) PinNext(ctx context.Context, fn link.FileName) (*Buffer, error) {
	b, err := bl.storage.RequestNext(ctx, fn)
	if err != nil {
		return nil, err
	}

	if bf, ok := bl.buffers[b.FileBlockID()]; ok {
		b.Release()
		return bf, nil
	}

	bf := NewBuffer(bl, b)
	bl.buffers[b.FileBlockID()] = bf
	return bf, nil
}

func (bl *BufferList) Release() {
	for _, b := range bl.buffers {
		b.b.Release()
	}
}

func (bl *BufferList) releaseBuffer(b *Buffer) {
	if !b.Dirty() {
		b.b.Release()
		delete(bl.buffers, b.FileBlockID())
	}
}
