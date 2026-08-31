package storage

import (
	"io"
	"sync"

	"github.com/liaradb/liaradb/storage/link"
)

type Buffer struct {
	blockID link.BlockID
	oldBID  link.BlockID
	data    []byte
	status  BufferStatus
	s       *Storage
	pins    int
	reads   uint
	mux     sync.RWMutex
	loader  func() error
}

func newBuffer(s *Storage) *Buffer {
	return &Buffer{
		data: make([]byte, s.BufferSize()),
		s:    s,
	}
}

func (b *Buffer) BlockID() link.BlockID { return b.blockID }
func (b *Buffer) Dirty() bool           { return b.status == BufferStatusDirty }
func (b *Buffer) Pins() int             { return b.pins }
func (b *Buffer) Reads() uint           { return b.reads }
func (b *Buffer) Size() int64           { return b.s.BufferSize() }
func (b *Buffer) Raw() []byte           { return b.data }
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
	b.reads++
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
func (b *Buffer) load(bid link.BlockID, next bool) {
	b.blockID = bid
	b.reads = 0
	b.initLoader(next)
}

// Move loading into sync.Once.
// This will allow loaded traffic to continue
func (b *Buffer) initLoader(
	next bool,
) {
	b.loader = sync.OnceValue(b.createLoader(next))
}

func (b *Buffer) createLoader(next bool) func() error {
	return func() error {
		if err := b.flushAndLoad(next); err != nil {
			b.initLoader(next)
			return err
		}

		return nil
	}
}

func (b *Buffer) flushAndLoad(next bool) error {
	if err := b.flushIfDirtyBeforeLoad(); err != nil {
		return err
	}

	// Only change oldBID after it has flushed
	b.oldBID = b.blockID
	b.status = BufferStatusLoading

	if err := b.clearOrLoad(next); err != nil {
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

func (b *Buffer) flushIfDirtyBeforeLoad() error {
	if !b.Dirty() {
		return nil
	}

	return b.flush()
}

func (b *Buffer) clearOrLoad(next bool) error {
	if next {
		clear(b.data)
		return nil
	}

	r, err := b.s.openFile(b.blockID)
	if err != nil {
		return err
	}

	if err := b.read(r); err != nil {
		b.status = BufferStatusCorrupt
		return err
	}

	return nil
}

func (b *Buffer) loadOnce() error {
	return b.loader()
}

func (b *Buffer) read(r io.ReaderAt) error {
	n, err := r.ReadAt(b.data, b.offset())
	if err != nil {
		// Ignore EOF
		if err != io.EOF {
			return err
		}

		// Clear the remainder of the buffer
		clear(b.data[n:])
	}

	return nil
}

func (b *Buffer) flushIfDirty() error {
	if !b.Dirty() {
		return nil
	}

	b.RLatch()
	defer b.RUnlatch()

	if err := b.flush(); err != nil {
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}

func (b *Buffer) flush() error {
	w, err := b.s.openFile(b.oldBID)
	if err != nil {
		return err
	}

	_, err = w.WriteAt(b.data, b.offsetOld())
	return err
}

func (b *Buffer) offset() int64 {
	return b.blockID.Offset(int64(len(b.data))).Value()
}

func (b *Buffer) offsetOld() int64 {
	return b.oldBID.Offset(int64(len(b.data))).Value()
}

func (b *Buffer) Clear() {
	clear(b.data)
	b.status = BufferStatusUninitialized
}

func (b *Buffer) Fill(data []byte) {
	n := copy(b.data, data)
	clear(b.data[n:])
	b.SetDirty()
}

func (b *Buffer) Clone(o *Buffer) {
	copy(b.Raw(), o.Raw())
	b.SetDirty()
}
