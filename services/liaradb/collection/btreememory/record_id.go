package btreememory

type RecordID struct {
	block    int64
	position int8
}

func NewRecordID(block int64, position int8) RecordID {
	return RecordID{
		block:    block,
		position: position,
	}
}

func (i RecordID) Block() int64   { return i.block }
func (i RecordID) Position() int8 { return i.position }
