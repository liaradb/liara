package storage

import "github.com/liaradb/liaradb/encoder/page"

type RecordID struct {
	BlockID  BlockID
	Position page.Offset
}

func (i RecordID) Offset(bufferSize int64) page.Offset {
	return i.BlockID.Offset(bufferSize) + i.Position
}
