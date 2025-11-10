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

// TODO: Test this
func (b BlockID) RecordID(position Offset) RecordID {
	return RecordID{
		BlockID:  b,
		Position: position,
	}
}
