package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/page"
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
	return i.blockID.Size() + i.slotID.Size()
}

func (i RecordID) Offset(bufferSize int64) page.Offset {
	return i.blockID.Offset(bufferSize) * page.Offset(i.slotID)
}

// TODO: Implement this
func (i *RecordID) Read(r io.Reader) error {
	return nil
}

// TODO: Implement this
func (i RecordID) Write(w io.Writer) error {
	return nil
}
