package storage

import (
	"encoding/binary"
	"errors"
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
	if pos < 0 || pos >= len(b.data)-8 {
		return errors.New("out of bounds")
	}

	binary.BigEndian.PutUint64(b.data[pos:], value)
	return nil
}

func (b *Buffer) ReadUint64(pos int) (uint64, error) {
	if pos < 0 || pos >= len(b.data)-8 {
		return 0, errors.New("out of bounds")
	}

	return binary.BigEndian.Uint64(b.data[pos:]), nil
}
