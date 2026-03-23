package storage

import (
	"io"
	"sync"

	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/storage/link"
)

type Buffer struct {
	blockID link.BlockID
	buffer  *buffer.Buffer
	status  BufferStatus
	s       *Storage
	pins    int
	mux     sync.RWMutex
}

func newBuffer(s *Storage) *Buffer {
	return &Buffer{
		buffer: buffer.New(s.BufferSize()),
		s:      s,
	}
}

func (b *Buffer) BlockID() link.BlockID { return b.blockID }
func (b *Buffer) Dirty() bool           { return b.status == BufferStatusDirty }
func (b *Buffer) Pins() int             { return b.pins }
func (b *Buffer) Size() int64           { return b.s.BufferSize() }
func (b *Buffer) Raw() []byte           { return b.buffer.Bytes() }
func (b *Buffer) Cursor() int64         { return b.buffer.Cursor() }
func (b *Buffer) Status() BufferStatus  { return b.status }

// TODO: Test these
func (b *Buffer) Latch()    { b.mux.Lock() }
func (b *Buffer) Unlatch()  { b.mux.Unlock() }
func (b *Buffer) RLatch()   { b.mux.RLock() }
func (b *Buffer) RUnlatch() { b.mux.RUnlock() }

// This is usually managed by the Buffer itself.
// However, it is useful when using Raw.
func (b *Buffer) SetDirty() { b.status = BufferStatusDirty }

func (b *Buffer) pin() {
	b.pins++
}

func (b *Buffer) unpin() bool {
	b.pins--
	if b.pins < 0 {
		// This should never happen
		panic("nevative pins")
	}
	return b.pins == 0
}

func (b *Buffer) Release() {
	b.s.release(b)
}

func (b *Buffer) load(bid link.BlockID, next bool) error {
	// blockID will always be changing
	// status is dirty only if already loaded
	if b.status == BufferStatusDirty {
		if err := b.s.flush(b); err != nil {
			return err
		}

		b.status = BufferStatusLoaded
	}

	b.blockID = bid
	b.status = BufferStatusLoading

	if next {
		// TODO: Test this
		b.buffer.Clear()
		b.status = BufferStatusLoaded
		return nil
	}

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
		// Ignore EOF
		if err != io.EOF {
			return err
		}

		// Clear the remainder of the buffer
		b.buffer.ClearAfter(n)
	}

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

// TODO: Test this
func (b *Buffer) flushIfDirty() error {
	if !b.Dirty() {
		return nil
	}

	b.RLatch()
	defer b.RUnlatch()

	if err := b.s.flush(b); err != nil {
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

func (b *Buffer) Clear() {
	b.buffer.Clear()
	b.status = BufferStatusUninitialized
}

func (b *Buffer) Read(p []byte) (int, error) {
	return b.buffer.Read(p)
}

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	return b.buffer.Seek(offset, whence)
}

func (b *Buffer) Write(p []byte) (int, error) {
	n, err := b.buffer.Write(p)
	if n != 0 {
		b.SetDirty()
	}

	return n, err
}

func (b *Buffer) Clone(o *Buffer) {
	copy(b.Raw(), o.Raw())
	b.SetDirty()
}
