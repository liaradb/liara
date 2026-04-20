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
	loader  func() error
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

// Load from file system
//   - blockID will always be changing
//   - status is dirty only if already loaded
func (b *Buffer) load(bid link.BlockID, next bool) error {
	w, err := b.s.openFile(b.blockID)
	if err != nil {
		return err
	}

	r, err := b.s.openFile(bid)
	if err != nil {
		return err
	}

	if err := b.flushIfDirtyBeforeLoad(w); err != nil {
		return err
	}

	b.blockID = bid

	// TODO: What do we do if this has an error?
	// Move loading into sync.Once.
	// This will allow loaded traffic to continue
	b.createLoad(func() error {
		b.status = BufferStatusLoading

		if err := b.clearOrLoad(next, r); err != nil {
			return err
		}

		b.status = BufferStatusLoaded
		return nil
	})

	return nil
}

func (b *Buffer) flushIfDirtyBeforeLoad(w io.WriterAt) error {
	if !b.Dirty() {
		return nil
	}

	return b.write(w)
}

func (b *Buffer) clearOrLoad(next bool, r io.ReaderAt) error {
	if next {
		b.buffer.Clear()
		return nil
	}

	if err := b.read(r); err != nil {
		b.status = BufferStatusCorrupt
		return err
	}

	return nil
}

func (b *Buffer) createLoad(once func() error) {
	b.loader = sync.OnceValue(once)
}

func (b *Buffer) loadOnce() error {
	return b.loader()
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
