package link

type BlockID struct {
	fileName FileName
	position FilePosition
}

func (b BlockID) FileName() FileName     { return b.fileName }
func (b BlockID) Position() FilePosition { return b.position }

func NewBlockID(fn FileName, position FilePosition) BlockID {
	return BlockID{
		fileName: fn,
		position: position,
	}
}

func (b BlockID) Size() int {
	return b.position.Size() + b.fileName.Size()
}

func (b BlockID) Offset(bufferSize int64) Offset {
	return b.position.Offset(bufferSize)
}

func (b *BlockID) SetPosition(p FilePosition) {
	b.position = p
}

func (b BlockID) RecordID(slotID SlotID) RecordID {
	return RecordID{
		blockID: b,
		slotID:  slotID,
	}
}

func (b BlockID) Next() BlockID {
	return BlockID{
		fileName: b.fileName,
		position: b.position + 1,
	}
}
