package storage

import "github.com/liaradb/liaradb/encoder/page"

type RecordID struct {
	BlockID  BlockID
	Position RecordPosition
}

func (i RecordID) Offset(bufferSize int64) page.Offset {
	return i.BlockID.Offset(bufferSize) + page.Offset(i.Position)
}
