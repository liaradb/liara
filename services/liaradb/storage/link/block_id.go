package link

import "github.com/liaradb/liaradb/encoder/page"

type BlockID struct {
	fileName FileName
	position FilePosition
}

func (b BlockID) FileName() FileName     { return b.fileName }
func (b BlockID) Position() FilePosition { return b.position }

func NewBlockID(fn FileName, position FilePosition) BlockID {
	return BlockID{
		fileName: fn,
		position: position,
	}
}

func (b BlockID) Offset(bufferSize int64) page.Offset {
	return b.position.Offset(bufferSize)
}

func (b *BlockID) SetPosition(p FilePosition) {
	b.position = p
}

func (b BlockID) RecordID(position RecordPosition) RecordID {
	return RecordID{
		blockID:  b,
		position: position,
	}
}
