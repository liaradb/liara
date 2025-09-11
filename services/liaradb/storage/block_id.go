package storage

import "github.com/liaradb/liaradb/raw"

type BlockID struct {
	FileName string
	Position raw.Offset
}

func (b BlockID) Offset(bufferSize int64) raw.Offset {
	return b.Position * raw.Offset(bufferSize)
}
