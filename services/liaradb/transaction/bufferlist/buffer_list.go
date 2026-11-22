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
	if bc, ok := bl.buffers[fbid]; ok {
		return bc.buffer, nil
	}

	b, err := bl.storage.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	bl.buffers[fbid] = &bufferCount{
		buffer: b,
		count:  1,
	}
	return b, nil
}

func (bl *BufferList) PinCurrent(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	b, err := bl.storage.RequestCurrent(ctx, fn)
	if err != nil {
		return nil, err
	}

	if bc, ok := bl.buffers[b.FileBlockID()]; ok {
		bc.count++
		return bc.buffer, nil
	}

	bl.buffers[b.FileBlockID()] = &bufferCount{
		buffer: b,
		count:  1,
	}
	return b, nil
}

func (bl *BufferList) PinNext(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	b, err := bl.storage.RequestNext(ctx, fn)
	if err != nil {
		return nil, err
	}

	if bc, ok := bl.buffers[b.FileBlockID()]; ok {
		bc.count++
		return bc.buffer, nil
	}

	bl.buffers[b.FileBlockID()] = &bufferCount{
		buffer: b,
		count:  1,
	}
	return b, nil
}

func (bl *BufferList) Release() {
	for _, bc := range bl.buffers {
		for range bc.count {
			bc.buffer.Release()
		}
	}
}

func (bl *BufferList) ReleaseBuffer(b *storage.Buffer) {
	if !b.Dirty() {
		bc, ok := bl.buffers[b.FileBlockID()]
		if !ok {
			panic("untracked buffer")
		}

		bc.buffer.Release()
		bc.count++
		if bc.count == 0 {
			delete(bl.buffers, b.FileBlockID())
		}
	}
}
