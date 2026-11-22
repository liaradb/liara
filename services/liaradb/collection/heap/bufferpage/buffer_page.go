package bufferpage

import (
	"github.com/liaradb/liaradb/collection/heap/record"
	"github.com/liaradb/liaradb/collection/log"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

// TODO: Track dirty pages inside and outside transaction
type BufferPage struct {
	page   record.Page
	buffer *storage.Buffer
	l      log.Log
}

func New(b *storage.Buffer, l log.Log) *BufferPage {
	return &BufferPage{
		page:   record.NewPage(b.Raw()),
		buffer: b,
		l:      l,
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
	bp.l.Unpin(bp.buffer)
}

func (bp *BufferPage) CommitAppend(size int16) {
	_ = bp.page.CommitAppend(size)
	// TODO: Verify this the correct time to do this
	bp.buffer.SetDirty()
}

func (bp *BufferPage) CommitUpdate() {
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

func (bp *BufferPage) Delete(i link.SlotID) {
	bp.page.Delete(i)
}
