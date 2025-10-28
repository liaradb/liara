package storage

import (
	"io"
	"iter"
	"sync"

	"github.com/liaradb/liaradb/page"
	"github.com/liaradb/liaradb/raw"
)

type Buffer struct {
	blockID BlockID
	buffer  *raw.Buffer
	page    *page.BytePage
	status  BufferStatus
	s       *Storage
	pins    int
	mux     sync.RWMutex
}

type BufferStatus int

const (
	BufferStatusUninitialized BufferStatus = iota
	BufferStatusLoading
	BufferStatusLoaded
	BufferStatusDirty
	BufferStatusCorrupt
)

func newBuffer(s *Storage) *Buffer {
	return &Buffer{
		buffer: raw.NewBuffer(s.BufferSize()),
		page:   page.New(page.Offset(s.BufferSize())),
		s:      s,
	}
}

func (b *Buffer) BlockID() BlockID { return b.blockID }
func (b *Buffer) Dirty() bool      { return b.status == BufferStatusDirty }
func (b *Buffer) Pins() int        { return b.pins }

// TODO: Test these
func (b *Buffer) Latch()    { b.mux.Lock() }
func (b *Buffer) Unlatch()  { b.mux.Unlock() }
func (b *Buffer) RLatch()   { b.mux.RLock() }
func (b *Buffer) RUnlatch() { b.mux.RUnlock() }

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

// TODO: Do we need to clone the item?
func (b *Buffer) Add(i []byte) error {
	b.Latch()
	defer b.Unlatch()

	return b.add(i)
}

func (b *Buffer) add(i []byte) error {
	if err := b.page.Add(page.NewItem(i)); err != nil {
		return err
	}

	b.status = BufferStatusDirty
	return nil
}

func (b *Buffer) Items() iter.Seq2[[]byte, error] {
	b.RLatch()
	defer b.RUnlatch()

	// TODO: Is there a simpler way?
	return func(yield func([]byte, error) bool) {
		for i, err := range b.page.Items() {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(i.Value(), nil) {
				return
			}
		}
	}
}

func (b *Buffer) read(r io.ReaderAt) error {
	n, err := r.ReadAt(b.buffer.Bytes(), b.offset())
	if err != nil {
		if err != io.EOF {
			return err
		}
		clear(b.buffer.Bytes()[n:])
	}

	// TODO: Test this
	if _, err := b.buffer.Seek(0, io.SeekStart); err != nil {
		return err
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
func (b *Buffer) FlushIfDirty() error {
	if !b.Dirty() {
		return nil
	}

	if err := b.s.flush(b); err != nil {
		return err
	}

	b.status = BufferStatusLoaded
	return nil
}
