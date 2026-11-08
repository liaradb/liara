package span

import (
	"context"
	"slices"

	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	fragmentHeaderSize = 0 +
		base.Uint64Size +
		base.Uint16Size +
		page.CrcSize
)

type Fragment struct {
	p           BufferPage
	slotID      link.SlotID
	nextBlockID wrap.Int64
	nextSlotID  wrap.Int16
	crc         wrap.Int32
	buffer      *buffer.Buffer
	size        int16
	base        []byte
}

type BufferPage interface {
	SetLogSequenceNumber(logpage.LogSequenceNumber)
	FileBlockID() link.FileBlockID
	CommitAppend(size int16)
	Release()
	Slot(link.SlotID) ([]byte, []byte)
}

func newFragment(
	p BufferPage,
	sid link.SlotID,
	header []byte,
	data []byte,
) *Fragment {
	nextBlockID, header0 := wrap.NewInt64(header)
	nextSlotID, header1 := wrap.NewInt16(header0)
	crc, _ := wrap.NewInt32(header1)
	return &Fragment{
		p:           p,
		slotID:      sid,
		nextBlockID: nextBlockID,
		nextSlotID:  nextSlotID,
		crc:         crc,
		buffer:      buffer.NewFromSlice(data),
		size:        int16(len(data)),
	}
}

func (f Fragment) length() int               { return f.buffer.Length() }
func (f Fragment) NextBlockID() link.BlockID { return link.BlockID(f.nextBlockID.Get()) }
func (f Fragment) NextSlotID() link.SlotID   { return link.SlotID(f.nextSlotID.Get()) }
func (f Fragment) RecordID() link.RecordID   { return f.p.FileBlockID().RecordID(f.slotID) }

func (f Fragment) valid() bool {
	return page.RestoreCRC(f.crc.GetUnsigned()).
		Compare(f.buffer.Bytes())
}

func (f Fragment) commit() {
	f.commitCRC()
	f.p.CommitAppend(f.size)
}

func (f Fragment) commitCRC() {
	crc := page.NewCRC(f.buffer.Bytes())
	f.crc.SetUnsigned(crc.Value())
}

func (f Fragment) setNextBlockID(bid link.BlockID) {
	f.nextBlockID.Set(bid.Value())
}

func (f Fragment) setNextSlotID(i link.SlotID) {
	f.nextSlotID.Set(i.Value())
}

func (f Fragment) Read(p []byte) (int, error) {
	return f.buffer.Read(p)
}

func (f Fragment) Write(p []byte) (int, error) {
	return f.buffer.Write(p)
}

func (f *Fragment) prepareWrite() {
	f.base = f.cloneBuffer()
}

func (f *Fragment) cloneBuffer() []byte {
	return slices.Clone(f.buffer.Bytes())
}

func (f *Fragment) logAppend(ctx context.Context, l Log) error {
	lsn, err := l.Append(ctx, f.fileRecordID(), f.cloneBuffer())
	if err != nil {
		return err
	}

	f.p.SetLogSequenceNumber(lsn)
	return nil
}

func (f *Fragment) logUpdate(ctx context.Context, l Log) error {
	lsn, err := l.Update(ctx, f.fileRecordID(), f.cloneBuffer(), f.base)
	if err != nil {
		return err
	}

	f.p.SetLogSequenceNumber(lsn)
	return nil
}

func (f Fragment) fileRecordID() link.FileRecordID {
	return f.p.FileBlockID().FileRecordID(f.slotID)
}
