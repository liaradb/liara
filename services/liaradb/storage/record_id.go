package storage

type RecordID struct {
	BlockID  BlockID
	Position Offset
}

func (i RecordID) Offset(bufferSize int64) Offset {
	return i.BlockID.Offset(bufferSize) + i.Position
}
