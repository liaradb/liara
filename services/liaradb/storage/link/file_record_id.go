package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

type FileRecordID struct {
	fileBlockID FileBlockID
	slotID      SlotID
}

func NewFileRecordID(
	fileBlockID FileBlockID,
	slotID SlotID,
) FileRecordID {
	return FileRecordID{
		fileBlockID: fileBlockID,
		slotID:      slotID,
	}
}

func (i FileRecordID) FileBlockID() FileBlockID { return i.fileBlockID }
func (i FileRecordID) SlotID() SlotID           { return i.slotID }

func (i FileRecordID) Size() int {
	return serializer.Size(
		i.fileBlockID,
		i.slotID)
}

func (i FileRecordID) Offset(bufferSize int64) Offset {
	return i.fileBlockID.Offset(bufferSize) * Offset(i.slotID)
}

func (i *FileRecordID) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&i.fileBlockID,
		&i.slotID)
}

func (i FileRecordID) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		i.fileBlockID,
		i.slotID)
}
