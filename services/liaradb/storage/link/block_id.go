package link

import "github.com/liaradb/liaradb/encoder/page"

type BlockID struct {
	fileName string
	position page.Offset
}

func (b BlockID) FileName() string      { return b.fileName }
func (b BlockID) Position() page.Offset { return b.position }

func NewBlockID(fileName string, position page.Offset) BlockID {
	return BlockID{
		fileName: fileName,
		position: position,
	}
}

func (b BlockID) Offset(bufferSize int64) page.Offset {
	return b.position * page.Offset(bufferSize)
}

func (b *BlockID) SetPosition(p page.Offset) {
	b.position = p
}

// TODO: Test this
func (b BlockID) RecordID(position RecordPosition) RecordID {
	return RecordID{
		blockID:  b,
		position: position,
	}
}
