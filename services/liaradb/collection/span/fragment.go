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
	FragmentHeaderSize = 0 +
		base.Uint64Size +
		base.Uint16Size +
		page.CrcSize
)

type Fragment struct {
	p            BufferPage
	slotID       link.SlotID
	nextPosition wrap.Int64
	nextSlotID   wrap.Int16
	crc          wrap.Int32
	buffer       *buffer.Buffer
	size         int
	base         []byte
}

type BufferPage interface {
	SetLogSequenceNumber(logpage.LogSequenceNumber)
	BlockID() link.BlockID
	Commit(size int)
	Release()
}

func newFragment(
	p BufferPage,
	sid link.SlotID,
	header []byte,
	data []byte,
) *Fragment {
	nextPosition, header0 := wrap.NewInt64(header)
	nextSlotID, header1 := wrap.NewInt16(header0)
	crc, _ := wrap.NewInt32(header1)
	return &Fragment{
		p:            p,
		slotID:       sid,
		nextPosition: nextPosition,
		nextSlotID:   nextSlotID,
		crc:          crc,
		buffer:       buffer.NewFromSlice(data),
		size:         len(data),
	}
}

// TODO: Fix this cast
func (f Fragment) length() int                     { return int(f.buffer.Length()) }
func (f Fragment) SlotID() link.SlotID             { return f.slotID } // TODO: This is unused
func (f Fragment) NextPosition() link.FilePosition { return link.FilePosition(f.nextPosition.Get()) }
func (f Fragment) NextSlotID() link.SlotID         { return link.SlotID(f.nextSlotID.Get()) }

// TODO: This is unused
func (f Fragment) BlockID(fn link.FileName) link.BlockID {
	return link.NewBlockID(fn, f.NextPosition())
}

// TODO: This is unused
func (f Fragment) RecordID(fn link.FileName) link.RecordID {
	return link.NewRecordID(f.BlockID(fn), f.NextSlotID())
}

func (f Fragment) valid() bool {
	return page.RestoreCRC(f.crc.GetUnsigned()).
		Compare(f.buffer.Bytes())
}

func (f Fragment) commit() {
	f.commitCRC()
	f.p.Commit(f.size)
}

func (f Fragment) commitCRC() {
	crc := page.NewCRC(f.buffer.Bytes())
	f.crc.SetUnsigned(crc.Value())
}

func (f Fragment) setNextPosition(p link.FilePosition) {
	f.nextPosition.Set(p.Value())
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

// TODO: Can we do this without the allocations?
func (f *Fragment) logAppend(ctx context.Context, l Log) error {
	lsn, err := l.Append(ctx, f.recordLocator(), f.cloneBuffer())
	if err != nil {
		return err
	}

	f.p.SetLogSequenceNumber(lsn)
	return nil
}

// TODO: Can we do this without the allocations?
func (f *Fragment) logUpdate(ctx context.Context, l Log) error {
	lsn, err := l.Update(ctx, f.recordLocator(), f.cloneBuffer(), f.base)
	if err != nil {
		return err
	}

	f.p.SetLogSequenceNumber(lsn)
	return nil
}

func (f Fragment) recordLocator() link.RecordLocator {
	return f.p.BlockID().RecordLocator(f.slotID)
}
