package bufferlist

import (
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type Buffer struct {
	bl    *BufferList
	b     *storage.Buffer
	dirty bool
}

func NewBuffer(bl *BufferList, b *storage.Buffer) *Buffer {
	return &Buffer{
		bl: bl,
		b:  b,
	}
}

func (b *Buffer) Raw() []byte                   { return b.b.Raw() }
func (b *Buffer) FileBlockID() link.FileBlockID { return b.b.FileBlockID() }
func (b *Buffer) Latch()                        { b.b.Latch() }
func (b *Buffer) Unlatch()                      { b.b.Unlatch() }
func (b *Buffer) RLatch()                       { b.b.RLatch() }
func (b *Buffer) RUnlatch()                     { b.b.RUnlatch() }

// Testing only
func (b *Buffer) Pins() int   { return b.b.Pins() }
func (b *Buffer) Reads() uint { return b.b.Reads() }

func (b *Buffer) Clone(a *Buffer) { b.b.Clone(a.b) }
func (b *Buffer) Release()        { b.bl.releaseBuffer(b) }

func (b *Buffer) Dirty() bool {
	return b.b.Dirty()
	// return b.dirty
}

func (b *Buffer) SetDirty() {
	b.b.SetDirty()
	b.dirty = true
}
