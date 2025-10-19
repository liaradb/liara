package storage

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/raw"
	"github.com/liaradb/liaradb/storage/record"
)

type Buffer struct {
	blockID BlockID
	buffer  *raw.Buffer
	page    *record.Page
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
		// TODO: This accesses BufferManager.bufferSize direclty
		buffer: raw.NewBuffer(s.bm.bufferSize),
		page:   record.NewPage(record.Offset(s.bm.bufferSize)),
		s:      s,
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

func (b *Buffer) Add(i []byte) error {
	if err := b.page.Add(i); err != nil {
		return err
	}

	b.status = BufferStatusDirty
	return nil
}

func (b *Buffer) Items() iter.Seq2[record.Item, error] {
	return b.page.Items()
}

func (b *Buffer) read(r io.ReaderAt) error {
	n, err := r.ReadAt(b.buffer.Bytes(), b.offset())
	if err != nil {
		if err != io.EOF {
			return err
		}
		clear(b.buffer.Bytes()[n:])
	}

	return b.page.Read(b.buffer)
}

func (b *Buffer) write(w io.WriterAt) error {
	b.buffer.Clear()
	if err := b.page.Write(b.buffer); err != nil {
		return err
	}

	_, err := w.WriteAt(b.buffer.Bytes(), b.offset())
	return err
}

func (b *Buffer) offset() int64 {
	return b.blockID.Offset(b.buffer.Length()).Value()
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
