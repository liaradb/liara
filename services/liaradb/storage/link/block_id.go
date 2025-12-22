package link

import (
	"github.com/liaradb/liaradb/encoder/page"
)

type BlockID struct {
	fileName FileName
	position page.Offset
}

func (b BlockID) FileName() FileName    { return b.fileName }
func (b BlockID) Position() page.Offset { return b.position }

func NewBlockID(fn FileName, position page.Offset) BlockID {
	return BlockID{
		fileName: fn,
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
