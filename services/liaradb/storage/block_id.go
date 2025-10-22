package storage

type BlockID struct {
	FileName string
	Position Offset
}

func NewBlockID(fileName string, position Offset) BlockID {
	return BlockID{
		FileName: fileName,
		Position: position,
	}
}

func (b BlockID) Offset(bufferSize int64) Offset {
	return b.Position * Offset(bufferSize)
}
