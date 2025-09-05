package storage

import "github.com/cardboardrobots/liaradb/raw"

type Buffer struct {
	blockID BlockID
	data    []byte
	dirty   bool
	bm      *BufferManager
}

func newBuffer(bid BlockID, bm *BufferManager) *Buffer {
	return &Buffer{
		blockID: bid,
		data:    make([]byte, bm.bufferSize),
		bm:      bm,
	}
}

func (b *Buffer) Dirty() bool { return b.dirty }

func (b *Buffer) Load() error {
	if err := b.bm.Load(b); err != nil {
		return err
	}

	b.dirty = false
	return nil
}

func (b *Buffer) Flush() error {
	if err := b.bm.Flush(b); err != nil {
		return err
	}

	b.dirty = false
	return nil
}

func (b *Buffer) WriteUint64(value uint64, off raw.Offset) error {
	if err := raw.CopyUint64(b.data, value, off); err != nil {
		return err
	}

	b.dirty = true
	return nil
}

func (b *Buffer) ReadUint64(off raw.Offset) (uint64, error) {
	return raw.GetUint64(b.data, off)
}
