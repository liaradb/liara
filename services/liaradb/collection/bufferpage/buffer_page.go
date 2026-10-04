package bufferpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type BufferPage struct {
	header
	page   *page.Page
	buffer *storage.Buffer
}

// TODO: Remove this parameter once the import cycle with span is fixed.
func New(b *storage.Buffer, slotHeaderSize int) *BufferPage {
	page := page.NewFromSlice(b.Raw(), headerSize, slotHeaderSize)
	header, _ := newHeader(page.Header())
	return &BufferPage{
		header: header,
		page:   page,
		buffer: b,
	}
}

func (bp *BufferPage) BlockID() link.BlockID { return bp.buffer.BlockID() }

func (bp *BufferPage) Shadow(base *BufferPage) {
	bp.page.Fill(base.page.Data())
}

func (bp *BufferPage) Release() {
	bp.buffer.Release()
}

func (bp *BufferPage) Commit(size int) {
	bp.page.Commit(size)
	// TODO: Verify this the correct time to do this
	bp.buffer.SetDirty()
}

func (bp *BufferPage) Count() link.SlotID {
	return bp.page.Count()
}

func (bp *BufferPage) Next(size int) ([]byte, []byte) {
	return bp.page.Next(size)
}

func (bp *BufferPage) Slot(i link.SlotID) ([]byte, []byte) {
	return bp.page.Slot(i)
}
