package storage

import "github.com/liaradb/liaradb/raw"

type BlockID struct {
	FileName string
	Position raw.Offset
}

func NewBlockID(fileName string, position raw.Offset) BlockID {
	return BlockID{
		FileName: fileName,
		Position: position,
	}
}

func (b BlockID) Offset(bufferSize int64) raw.Offset {
	return b.Position * raw.Offset(bufferSize)
}
