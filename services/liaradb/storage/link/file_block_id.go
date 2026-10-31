package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

type FileBlockID struct {
	fileName FileName
	position FilePosition
}

func (b FileBlockID) FileName() FileName     { return b.fileName }
func (b FileBlockID) Position() FilePosition { return b.position }

func NewFileBlockID(fn FileName, position FilePosition) FileBlockID {
	return FileBlockID{
		fileName: fn,
		position: position,
	}
}

func (b FileBlockID) Size() int {
	return serializer.Size(
		b.fileName,
		b.position)
}

func (b FileBlockID) Offset(bufferSize int64) Offset {
	return b.position.Offset(bufferSize)
}

func (b *FileBlockID) SetPosition(p FilePosition) {
	b.position = p
}

func (b FileBlockID) RecordID(slotID SlotID) RecordID {
	return RecordID{
		fileBlockID: b,
		slotID:      slotID,
	}
}

func (b FileBlockID) RecordLocator(slotID SlotID) RecordLocator {
	return NewRecordLocator(b.position, slotID)
}

func (b FileBlockID) Next() FileBlockID {
	return FileBlockID{
		fileName: b.fileName,
		position: b.position + 1,
	}
}

func (i *FileBlockID) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&i.fileName,
		&i.position)
}

func (i FileBlockID) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		i.fileName,
		i.position)
}
