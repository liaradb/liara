package storage

import "context"

type BufferList struct {
	storage *Storage
	buffers map[BlockID]*Buffer
}

func NewBufferList(storage *Storage) *BufferList {
	return &BufferList{
		storage: storage,
		buffers: make(map[BlockID]*Buffer),
	}
}

func (bl *BufferList) Pin(ctx context.Context, bid BlockID) (*Buffer, error) {
	if b, ok := bl.buffers[bid]; ok {
		return b, nil
	}

	b, err := bl.storage.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	bl.buffers[bid] = b
	return b, nil
}

func (bl *BufferList) Release() {
	for _, b := range bl.buffers {
		b.Release()
	}
}

func (bl *BufferList) ReleaseBuffer(b *Buffer) {
	b.Release()
	delete(bl.buffers, b.BlockID())
}
