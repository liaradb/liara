package storage

import "github.com/liaradb/liaradb/encoder/page"

type BlockID struct {
	FileName string
	Position page.Offset
}

func NewBlockID(fileName string, position page.Offset) BlockID {
	return BlockID{
		FileName: fileName,
		Position: position,
	}
}

func (b BlockID) Offset(bufferSize int64) page.Offset {
	return b.Position * page.Offset(bufferSize)
}

// TODO: Test this
func (b BlockID) RecordID(position RecordPosition) RecordID {
	return RecordID{
		blockID:  b,
		position: position,
	}
}
