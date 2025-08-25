package storage

import (
	"github.com/cardboardrobots/liaradb/raw"
)

type Position int64

type BlockID struct {
	FileName string
	Position Position
}

type Buffer struct {
	blockID BlockID
	data    []byte
	bm      *BufferManager
}

func (b *Buffer) Load() error {
	return b.bm.Load(b)
}

func (b *Buffer) Flush() error {
	return b.bm.Flush(b)
}

func (b *Buffer) WriteUint64(value uint64, pos int) error {
	return raw.CopyUint64(b.data, value, raw.Offset(pos))
}

func (b *Buffer) ReadUint64(pos int) (uint64, error) {
	return raw.GetUint64(b.data, raw.Offset(pos))
}
