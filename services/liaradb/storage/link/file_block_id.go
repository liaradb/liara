package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

type FileBlockID struct {
	fileName FileName
	blockID  BlockID
}

func (b FileBlockID) FileName() FileName { return b.fileName }
func (b FileBlockID) BlockID() BlockID   { return b.blockID }

func NewFileBlockID(fn FileName, blockID BlockID) FileBlockID {
	return FileBlockID{
		fileName: fn,
		blockID:  blockID,
	}
}

func (b FileBlockID) Size() int {
	return serializer.Size(
		b.fileName,
		b.blockID)
}

func (b FileBlockID) Offset(bufferSize int64) Offset {
	return b.blockID.Offset(bufferSize)
}

func (b *FileBlockID) SetBlockID(p BlockID) {
	b.blockID = p
}

func (b FileBlockID) RecordID(slotID SlotID) RecordID {
	return RecordID{
		fileBlockID: b,
		slotID:      slotID,
	}
}

func (b FileBlockID) RecordLocator(slotID SlotID) RecordLocator {
	return NewRecordLocator(b.blockID, slotID)
}

func (b FileBlockID) Next() FileBlockID {
	return FileBlockID{
		fileName: b.fileName,
		blockID:  b.blockID + 1,
	}
}

func (i *FileBlockID) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&i.fileName,
		&i.blockID)
}

func (i FileBlockID) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		i.fileName,
		i.blockID)
}
