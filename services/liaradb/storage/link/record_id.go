package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

type RecordID struct {
	fileBlockID FileBlockID
	slotID      SlotID
}

func NewRecordID(
	fileBlockID FileBlockID,
	slotID SlotID,
) RecordID {
	return RecordID{
		fileBlockID: fileBlockID,
		slotID:      slotID,
	}
}

func (i RecordID) FileBlockID() FileBlockID { return i.fileBlockID }
func (i RecordID) SlotID() SlotID           { return i.slotID }

func (i RecordID) Size() int {
	return serializer.Size(
		i.fileBlockID,
		i.slotID)
}

func (i RecordID) Offset(bufferSize int64) Offset {
	return i.fileBlockID.Offset(bufferSize) * Offset(i.slotID)
}

func (i *RecordID) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&i.fileBlockID,
		&i.slotID)
}

func (i RecordID) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		i.fileBlockID,
		i.slotID)
}
