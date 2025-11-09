package eventlog

import (
	"io"
	"iter"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
)

type BufferPage struct {
	buffer *storage.Buffer
	page   *page.BytePage
}

func NewBufferPage(b *storage.Buffer) *BufferPage {
	return &BufferPage{
		buffer: b,
		page:   page.New(page.Offset(b.Size())), // TODO: Find a better way to get size
	}
}

func (b *BufferPage) BlockID() storage.BlockID { return b.buffer.BlockID() }
func (b *BufferPage) Dirty() bool              { return b.buffer.Dirty() }
func (b *BufferPage) Pins() int                { return b.buffer.Pins() }
func (b *BufferPage) Flush() error             { return b.buffer.Flush() }
func (b *BufferPage) Release()                 { b.buffer.Release() }

// TODO: Do we need to clone the item?
func (b *BufferPage) Add(i []byte) error {
	b.buffer.Latch()
	defer b.buffer.Unlatch()

	if err := b.Unpack(); err != nil {
		return err
	}

	if err := b.add(i); err != nil {
		return err
	}

	return b.Pack()
}

func (b *BufferPage) add(i []byte) error {
	if err := b.page.Add(page.NewItem(i)); err != nil {
		return err
	}

	return nil
}

func (b *BufferPage) Items() iter.Seq2[[]byte, error] {
	b.buffer.RLatch()
	defer b.buffer.RUnlatch()

	// TODO: Is there a simpler way?
	return func(yield func([]byte, error) bool) {
		if err := b.Unpack(); err != nil {
			yield(nil, err)
			return
		}

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

// TODO: Don't access buffer.buffer directly
func (b *BufferPage) Unpack() error {
	if _, err := b.buffer.Seek(0, io.SeekStart); err != nil {
		return err
	}

	return b.page.Read(b.buffer)
}

// TODO: Don't access buffer.buffer directly
func (b *BufferPage) Pack() error {
	b.buffer.Clear()
	if err := b.page.Write(b.buffer); err != nil {
		return err
	}

	return nil
}
