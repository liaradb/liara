package link

import "github.com/liaradb/liaradb/encoder/page"

type RecordID struct {
	blockID BlockID
	slotID  SlotID
}

func NewRecordID(
	blockID BlockID,
	slotID SlotID,
) RecordID {
	return RecordID{
		blockID: blockID,
		slotID:  slotID,
	}
}

func (i RecordID) BlockID() BlockID { return i.blockID }
func (i RecordID) SlotID() SlotID   { return i.slotID }

func (i RecordID) Offset(bufferSize int64) page.Offset {
	return i.blockID.Offset(bufferSize) * page.Offset(i.slotID)
}
