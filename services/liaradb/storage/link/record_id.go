package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

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

func (i RecordID) Size() int {
	return serializer.Size(
		i.blockID,
		i.slotID)
}

func (i RecordID) Offset(bufferSize int64) Offset {
	return i.blockID.Offset(bufferSize) * Offset(i.slotID)
}

func (i *RecordID) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&i.blockID,
		&i.slotID)
}

func (i RecordID) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		i.blockID,
		i.slotID)
}
