package storage

import "github.com/cardboardrobots/liaradb/raw"

type Buffer struct {
	blockID BlockID
	data    []byte
	bm      *BufferManager
}

func newBuffer(bid BlockID, bm *BufferManager) *Buffer {
	return &Buffer{
		blockID: bid,
		data:    make([]byte, bm.bufferSize),
		bm:      bm,
	}
}

func (b *Buffer) Load() error {
	return b.bm.Load(b)
}

func (b *Buffer) Flush() error {
	return b.bm.Flush(b)
}

func (b *Buffer) WriteUint64(value uint64, off raw.Offset) error {
	return raw.CopyUint64(b.data, value, off)
}

func (b *Buffer) ReadUint64(off raw.Offset) (uint64, error) {
	return raw.GetUint64(b.data, off)
}
