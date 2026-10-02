package bufferpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type BufferPage struct {
	*page.Page
	header
	buffer *storage.Buffer
}

// TODO: Remove this parameter once the import cycle with span is fixed.
func New(b *storage.Buffer, slotHeaderSize int) *BufferPage {
	page := page.NewFromSlice(b.Raw(), headerSize, slotHeaderSize)
	header, _ := newHeader(page.Header())
	return &BufferPage{
		Page:   page,
		header: header,
		buffer: b,
	}
}

func (bp *BufferPage) BlockID() link.BlockID { return bp.buffer.BlockID() }

func (bp *BufferPage) Fill(data []byte) {
	bp.Page.Fill(data)
}

func (bp *BufferPage) Shadow(base *BufferPage) {
	bp.Page.Fill(base.Data())
}

func (bp *BufferPage) Release() {
	bp.buffer.Release()
}

func (bp *BufferPage) Commit(size int) {
	bp.Page.Commit(size)
	// TODO: Verify this the correct time to do this
	bp.buffer.SetDirty()
}
