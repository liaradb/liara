package link

import "github.com/liaradb/liaradb/encoder/page"

type RecordID struct {
	blockID  BlockID
	position RecordPosition
}

func NewRecordID(
	blockID BlockID,
	position RecordPosition,
) RecordID {
	return RecordID{
		blockID:  blockID,
		position: position,
	}
}

func (i RecordID) BlockID() BlockID         { return i.blockID }
func (i RecordID) Position() RecordPosition { return i.position }

func (i RecordID) Offset(bufferSize int64) page.Offset {
	return i.blockID.Offset(bufferSize) * page.Offset(i.position)
}
