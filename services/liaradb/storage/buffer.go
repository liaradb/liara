package storage

import "github.com/cardboardrobots/liaradb/raw"

type Buffer struct {
	blockID BlockID
	data    []byte
	status  BufferStatus
	bm      *BufferManager
	pins    int
}

type BufferStatus int

const (
	BufferStatusUninitialized BufferStatus = iota
	BufferStatusLoading
	BufferStatusLoaded
	BufferStatusDirty
	BufferStatusCorrupt
)

func newBuffer(bm *BufferManager) *Buffer {
	return &Buffer{
		data: make([]byte, bm.bufferSize),
		bm:   bm,
	}
}

func (b *Buffer) Dirty() bool { return b.status == BufferStatusDirty }
func (b *Buffer) Pins() int   { return b.pins }

func (b *Buffer) pin() {
	b.pins++
}

func (b *Buffer) unpin() bool {
	b.pins--
	// TODO: Do we need this?
	if b.pins < 0 {
		b.pins = 0
	}
	return b.pins == 0
}

func (b *Buffer) Load(bid BlockID) error {
	if b.blockID != bid && b.status == BufferStatusDirty {
		if err := b.bm.Flush(b); err != nil {
			return err
		}
	}

	b.blockID = bid
	b.status = BufferStatusLoading

	if err := b.bm.Load(b); err != nil {
		b.status = BufferStatusCorrupt
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

func (b *Buffer) Flush() error {
	if b.status != BufferStatusDirty {
		// TODO: Do we need more specific errors?
		return ErrNotDirty
	}

	if err := b.bm.Flush(b); err != nil {
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

func (b *Buffer) WriteUint64(value uint64, off raw.Offset) error {
	if err := raw.CopyUint64(b.data, value, off); err != nil {
		return err
	}

	b.status = BufferStatusDirty
	return nil
}

func (b *Buffer) ReadUint64(off raw.Offset) (uint64, error) {
	return raw.GetUint64(b.data, off)
}
