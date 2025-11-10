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
		page:   page.New(page.Offset(b.Size())),
	}
}

func (b *BufferPage) BlockID() storage.BlockID { return b.buffer.BlockID() }
func (b *BufferPage) Dirty() bool              { return b.buffer.Dirty() }
func (b *BufferPage) Pins() int                { return b.buffer.Pins() }
func (b *BufferPage) Flush() error             { return b.buffer.Flush() }
func (b *BufferPage) Release()                 { b.buffer.Release() }

func (b *BufferPage) Add(i []byte) (page.Offset, error) {
	b.buffer.Latch()
	defer b.buffer.Unlatch()

	if err := b.unpack(); err != nil {
		return 0, err
	}

	offset, err := b.add(i)
	if err != nil {
		return 0, err
	}

	if err := b.pack(); err != nil {
		return 0, err
	}

	return offset, nil
}

func (b *BufferPage) add(i []byte) (page.Offset, error) {
	return b.page.Add(page.NewItem(i))
}

func (b *BufferPage) Items() iter.Seq2[[]byte, error] {
	b.buffer.RLatch()
	defer b.buffer.RUnlatch()

	// TODO: Is there a simpler way?
	return func(yield func([]byte, error) bool) {
		if err := b.unpack(); err != nil {
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

func (b *BufferPage) unpack() error {
	if _, err := b.buffer.Seek(0, io.SeekStart); err != nil {
		return err
	}

	return b.page.Read(b.buffer)
}

func (b *BufferPage) pack() error {
	b.buffer.Clear()
	if err := b.page.Write(b.buffer); err != nil {
		return err
	}

	return nil
}
