package bufferpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	headerSize = 0
)

type BufferPage struct {
	page   *page.Page
	buffer *storage.Buffer
}

// TODO: Remove this parameter once the import cycle with span is fixed.
func New(b *storage.Buffer, slotHeaderSize int) *BufferPage {
	page := page.NewFromSlice(b.Raw(), headerSize, slotHeaderSize)
	return &BufferPage{
		page:   page,
		buffer: b,
	}
}

func (h *BufferPage) LogSequenceNumber() logpage.LogSequenceNumber {
	return logpage.NewLogSequenceNumber(h.page.TrackingID())
}

func (h *BufferPage) SetLogSequenceNumber(lsn logpage.LogSequenceNumber) {
	h.page.SetTrackingID(lsn.Value())
}

func (bp *BufferPage) FileBlockID() link.FileBlockID { return bp.buffer.FileBlockID() }

func (bp *BufferPage) Shadow(base *BufferPage) {
	bp.page.Fill(base.page.Data())
}

// TODO: Must latch before write
func (bp *BufferPage) Release() {
	bp.buffer.Release()
}

func (bp *BufferPage) CommitAppend(size int) {
	_ = bp.page.CommitAppend(size)
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
