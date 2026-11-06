package bufferpage

import (
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type BufferPage struct {
	page   span.Page
	buffer *storage.Buffer
}

func New(b *storage.Buffer) *BufferPage {
	return &BufferPage{
		page:   span.NewPage(b.Raw()),
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

func (bp *BufferPage) CommitAppend(size int16) {
	_ = bp.page.CommitAppend(size)
	// TODO: Verify this the correct time to do this
	bp.buffer.SetDirty()
}

func (bp *BufferPage) Count() link.SlotID {
	return bp.page.Count()
}

func (bp *BufferPage) Next(size int16) ([]byte, []byte) {
	return bp.page.Next(size)
}

func (bp *BufferPage) Slot(i link.SlotID) ([]byte, []byte) {
	return bp.page.Slot(i)
}
