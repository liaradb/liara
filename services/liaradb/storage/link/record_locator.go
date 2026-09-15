package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

const RecordLocatorSize = FilePositionSize + SlotIDSize

type RecordLocator struct {
	block  FilePosition
	slotID SlotID
}

func NewRecordLocator(block FilePosition, slotID SlotID) RecordLocator {
	return RecordLocator{
		block:  block,
		slotID: slotID,
	}
}

func (rl RecordLocator) Block() FilePosition { return rl.block }
func (rl RecordLocator) SlotID() SlotID      { return rl.slotID }
func (rl RecordLocator) Size() int           { return RecordLocatorSize }

func (rl RecordLocator) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		rl.block,
		rl.slotID)
}

func (rl *RecordLocator) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&rl.block,
		&rl.slotID)
}

func (rl RecordLocator) WriteData(data []byte) ([]byte, bool) {
	data0, ok := rl.block.WriteData(data)
	if !ok {
		return nil, false
	}

	return rl.slotID.WriteData(data0)
}

func (rl *RecordLocator) ReadData(data []byte) ([]byte, bool) {
	data0, ok := rl.block.ReadData(data)
	if !ok {
		return nil, false
	}

	return rl.slotID.ReadData(data0)
}
