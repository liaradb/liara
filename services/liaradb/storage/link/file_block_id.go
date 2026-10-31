package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

type FileBlockID struct {
	fileName FileName
	blockID  BlockID
}

func (i FileBlockID) FileName() FileName { return i.fileName }
func (i FileBlockID) BlockID() BlockID   { return i.blockID }

func NewFileBlockID(fn FileName, blockID BlockID) FileBlockID {
	return FileBlockID{
		fileName: fn,
		blockID:  blockID,
	}
}

func (i FileBlockID) Size() int {
	return serializer.Size(
		i.fileName,
		i.blockID)
}

func (i FileBlockID) Offset(bufferSize int64) Offset {
	return i.blockID.Offset(bufferSize)
}

func (i *FileBlockID) SetBlockID(p BlockID) {
	i.blockID = p
}

func (i FileBlockID) FileRecordID(slotID SlotID) FileRecordID {
	return FileRecordID{
		fileBlockID: i,
		slotID:      slotID,
	}
}

func (i FileBlockID) RecordID(slotID SlotID) RecordID {
	return NewRecordID(i.blockID, slotID)
}

func (i FileBlockID) Next() FileBlockID {
	return FileBlockID{
		fileName: i.fileName,
		blockID:  i.blockID + 1,
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
