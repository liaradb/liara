package record

import (
	"context"
	"slices"

	"github.com/liaradb/liaradb/collection/log"
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

type fragment struct {
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
	Delete(link.SlotID)
}

func newFragment(
	p BufferPage,
	sid link.SlotID,
	header []byte,
	data []byte,
) *fragment {
	nextBlockID, header0 := wrap.NewInt64(header)
	nextSlotID, header1 := wrap.NewInt16(header0)
	crc, _ := wrap.NewInt32(header1)
	return &fragment{
		p:           p,
		slotID:      sid,
		nextBlockID: nextBlockID,
		nextSlotID:  nextSlotID,
		crc:         crc,
		buffer:      buffer.NewFromSlice(data),
		size:        int16(len(data)),
	}
}

func (f fragment) length() int               { return f.buffer.Length() }
func (f fragment) NextBlockID() link.BlockID { return link.BlockID(f.nextBlockID.Get()) }
func (f fragment) NextSlotID() link.SlotID   { return link.SlotID(f.nextSlotID.Get()) }
func (f fragment) RecordID() link.RecordID   { return f.p.FileBlockID().RecordID(f.slotID) }

func (f fragment) valid() bool {
	return page.RestoreCRC(f.crc.GetUnsigned()).
		Compare(f.buffer.Bytes())
}

func (f fragment) commit() {
	f.commitCRC()
	f.p.CommitAppend(f.size)
}

func (f fragment) commitCRC() {
	crc := page.NewCRC(f.buffer.Bytes())
	f.crc.SetUnsigned(crc.Value())
}

func (f fragment) setNextBlockID(bid link.BlockID) {
	f.nextBlockID.Set(bid.Value())
}

func (f fragment) setNextSlotID(i link.SlotID) {
	f.nextSlotID.Set(i.Value())
}

func (f fragment) Read(p []byte) (int, error) {
	return f.buffer.Read(p)
}

func (f fragment) Write(p []byte) (int, error) {
	return f.buffer.Write(p)
}

func (f *fragment) prepareWrite() {
	f.base = f.cloneBuffer()
}

func (f *fragment) cloneBuffer() []byte {
	return slices.Clone(f.buffer.Bytes())
}

func (f *fragment) logAppend(ctx context.Context, l log.Log) error {
	lsn, err := l.Append(ctx, f.fileRecordID(), f.cloneBuffer())
	if err != nil {
		return err
	}

	f.p.SetLogSequenceNumber(lsn)
	return nil
}

func (f *fragment) logUpdate(ctx context.Context, l log.Log) error {
	lsn, err := l.Update(ctx, f.fileRecordID(), f.cloneBuffer(), f.base)
	if err != nil {
		return err
	}

	f.p.SetLogSequenceNumber(lsn)
	return nil
}

func (f *fragment) fileRecordID() link.FileRecordID {
	return f.p.FileBlockID().FileRecordID(f.slotID)
}

func (f *fragment) Delete() {
	f.p.Delete(f.slotID)
}
