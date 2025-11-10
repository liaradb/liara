package storage

import (
	"io"
	"sync"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Buffer struct {
	blockID BlockID
	buffer  *raw.Buffer
	status  BufferStatus
	s       *Storage
	pins    int
	mux     sync.RWMutex
}

func newBuffer(s *Storage) *Buffer {
	return &Buffer{
		buffer: raw.NewBuffer(s.BufferSize()),
		s:      s,
	}
}

func (b *Buffer) BlockID() BlockID { return b.blockID }
func (b *Buffer) Dirty() bool      { return b.status == BufferStatusDirty }
func (b *Buffer) Pins() int        { return b.pins }
func (b *Buffer) Size() int64      { return b.s.BufferSize() }

// TODO: Test these
func (b *Buffer) Latch()    { b.mux.Lock() }
func (b *Buffer) Unlatch()  { b.mux.Unlock() }
func (b *Buffer) RLatch()   { b.mux.RLock() }
func (b *Buffer) RUnlatch() { b.mux.RUnlock() }

func (b *Buffer) setDirty() { b.status = BufferStatusDirty }

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
func (b *Buffer) load(bid BlockID) error {
	if b.blockID != bid && b.status == BufferStatusDirty {
		if err := b.s.flush(b); err != nil {
			return err
		}
	}

	b.blockID = bid
	b.status = BufferStatusLoading

	if err := b.s.load(b); err != nil {
		b.status = BufferStatusCorrupt
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

func (b *Buffer) read(r io.ReaderAt) error {
	n, err := r.ReadAt(b.buffer.Bytes(), b.offset())
	if err != nil {
		if err != io.EOF {
			return err
		}
		clear(b.buffer.Bytes()[n:])
		// TODO: Test if page has been initialized
		// if n == 0 {
		// 	// TODO: Initialize page
		// 	b.page.Reset(page.ZeroHeader{})
		// 	return nil
		// }
	}

	// TODO: Test this
	_, err = b.buffer.Seek(0, io.SeekStart)
	return err
}

func (b *Buffer) write(w io.WriterAt) error {
	_, err := w.WriteAt(b.buffer.Bytes(), b.offset())
	return err
}

func (b *Buffer) offset() int64 {
	return b.blockID.Offset(b.buffer.Length()).Value()
}

// TODO: Do we need to latch?
func (b *Buffer) Flush() error {
	if !b.Dirty() {
		// TODO: Do we need more specific errors?
		return ErrNotDirty
	}

	if err := b.s.flush(b); err != nil {
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

// TODO: Test this
func (b *Buffer) flushIfDirty() error {
	if !b.Dirty() {
		return nil
	}

	if err := b.s.flush(b); err != nil {
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

// TODO: Test this
func (b *Buffer) Clear() {
	b.buffer.Clear()
	b.status = BufferStatusUninitialized
}

// TODO: Test this
func (b *Buffer) Read(p []byte) (int, error) {
	return b.buffer.Read(p)
}

// TODO: Test this
func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	return b.buffer.Seek(offset, whence)
}

// TODO: Test this
func (b *Buffer) Write(p []byte) (int, error) {
	n, err := b.buffer.Write(p)
	if n != 0 {
		b.setDirty()
	}

	return n, err
}
