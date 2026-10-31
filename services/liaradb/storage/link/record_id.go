package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

const RecordIDSize = BlockIDSize + SlotIDSize

type RecordID struct {
	block  BlockID
	slotID SlotID
}

func NewRecordID(block BlockID, slotID SlotID) RecordID {
	return RecordID{
		block:  block,
		slotID: slotID,
	}
}

func (i RecordID) Block() BlockID { return i.block }
func (i RecordID) SlotID() SlotID { return i.slotID }
func (i RecordID) Size() int      { return RecordIDSize }

func (i RecordID) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		i.block,
		i.slotID)
}

func (i *RecordID) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&i.block,
		&i.slotID)
}

func (i RecordID) WriteData(data []byte) []byte {
	data0 := i.block.WriteData(data)
	return i.slotID.WriteData(data0)
}

func (i *RecordID) ReadData(data []byte) []byte {
	data0 := i.block.ReadData(data)
	return i.slotID.ReadData(data0)
}
