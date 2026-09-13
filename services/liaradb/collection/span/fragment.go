package span

import (
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
	l            Log
	p            BufferPage
	slotID       link.SlotID
	nextPosition wrap.Int64
	nextSlotID   wrap.Int16
	crc          wrap.Int32
	buffer       *buffer.Buffer
}

type BufferPage interface {
	SetLogSequenceNumber(logpage.LogSequenceNumber)
	BlockID() link.BlockID
}

func newFragment(
	l Log,
	p BufferPage,
	sid link.SlotID,
	header []byte,
	data []byte,
) *Fragment {
	nextPosition, header0 := wrap.NewInt64(header)
	nextSlotID, header1 := wrap.NewInt16(header0)
	crc, _ := wrap.NewInt32(header1)
	return &Fragment{
		l:            l,
		p:            p,
		slotID:       sid,
		nextPosition: nextPosition,
		nextSlotID:   nextSlotID,
		crc:          crc,
		buffer:       buffer.NewFromSlice(data),
	}
}

// TODO: Fix this cast
func (f Fragment) length() int                     { return int(f.buffer.Length()) }
func (f Fragment) SlotID() link.SlotID             { return f.slotID }
func (f Fragment) NextPosition() link.FilePosition { return link.FilePosition(f.nextPosition.Get()) }
func (f Fragment) NextSlotID() link.SlotID         { return link.SlotID(f.nextSlotID.Get()) }

func (f Fragment) BlockID(fn link.FileName) link.BlockID {
	return link.NewBlockID(fn, f.NextPosition())
}

func (f Fragment) RecordID(fn link.FileName) link.RecordID {
	return link.NewRecordID(f.BlockID(fn), f.NextSlotID())
}

func (f Fragment) valid() bool {
	return page.RestoreCRC(f.crc.Get()).
		Compare(f.buffer.Bytes())
}

func (f Fragment) commit() {
	crc := page.NewCRC(f.buffer.Bytes())
	f.crc.Set(int32(crc.Value()))
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
