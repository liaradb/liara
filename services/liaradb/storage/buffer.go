package storage

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type Buffer struct {
	blockID BlockID
	data    []byte
	status  BufferStatus
	s       *Storage
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

// TODO: This should be private
func NewBuffer(s *Storage) *Buffer {
	return &Buffer{
		data: make([]byte, s.bm.bufferSize),
		s:    s,
	}
}

func (b *Buffer) BlockID() BlockID { return b.blockID }
func (b *Buffer) Dirty() bool      { return b.status == BufferStatusDirty }
func (b *Buffer) Pins() int        { return b.pins }

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

func (b *Buffer) Release() {
	b.s.release(b)
}

// TODO: Only load if BlockID is changing
func (b *Buffer) Load(bid BlockID) error {
	if b.blockID != bid && b.status == BufferStatusDirty {
		if err := b.s.bm.Flush(b); err != nil {
			return err
		}
	}

	b.blockID = bid
	b.status = BufferStatusLoading

	if err := b.s.bm.Load(b); err != nil {
		b.status = BufferStatusCorrupt
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

func (b *Buffer) read(r io.ReaderAt) error {
	_, err := r.ReadAt(b.data, b.offset())
	return err
}

func (b *Buffer) write(w io.WriterAt) error {
	_, err := w.WriteAt(b.data, b.offset())
	return err
}

func (b *Buffer) offset() int64 {
	return b.blockID.Offset(int64(len(b.data))).Value()
}

func (b *Buffer) Flush() error {
	if b.status != BufferStatusDirty {
		// TODO: Do we need more specific errors?
		return ErrNotDirty
	}

	if err := b.s.bm.Flush(b); err != nil {
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
