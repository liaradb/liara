package storage

import "github.com/liaradb/liaradb/raw"

type RecordID struct {
	BlockID  BlockID
	Position raw.Offset
}

func (i RecordID) Offset(bufferSize int64) raw.Offset {
	return i.BlockID.Offset(bufferSize) + i.Position
}
