package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

const RecordLocatorSize = BlockIDSize + SlotIDSize

type RecordLocator struct {
	block  BlockID
	slotID SlotID
}

func NewRecordLocator(block BlockID, slotID SlotID) RecordLocator {
	return RecordLocator{
		block:  block,
		slotID: slotID,
	}
}

func (rl RecordLocator) Block() BlockID { return rl.block }
func (rl RecordLocator) SlotID() SlotID { return rl.slotID }
func (rl RecordLocator) Size() int      { return RecordLocatorSize }

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

func (rl RecordLocator) WriteData(data []byte) []byte {
	data0 := rl.block.WriteData(data)
	return rl.slotID.WriteData(data0)
}

func (rl *RecordLocator) ReadData(data []byte) []byte {
	data0 := rl.block.ReadData(data)
	return rl.slotID.ReadData(data0)
}
