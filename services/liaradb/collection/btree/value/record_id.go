package value

import (
	"github.com/liaradb/liaradb/encoder/wrap"
)

const RecordIDSize = 8 + 1

// TODO: Test this
type RecordID struct {
	block    BlockPosition
	position RecordPosition
}

func NewRecordID(block BlockPosition, position RecordPosition) RecordID {
	return RecordID{
		block:    block,
		position: position,
	}
}

func (i RecordID) Block() int64   { return i.block.Value() }
func (i RecordID) Position() int8 { return i.position.Value() }

func (i RecordID) Size() int { return RecordIDSize }

// TODO: Use simpler wrap
func (le RecordID) Write(data []byte) []byte {
	block, data0 := wrap.NewInt64(data)
	position, data1 := wrap.NewByte(data0)

	block.Set(le.block.Value())
	position.Set(le.position.Value())

	return data1
}

// TODO: Use simpler wrap
func (le *RecordID) Read(data []byte) []byte {
	block, data0 := wrap.NewInt64(data)
	position, data1 := wrap.NewByte(data0)

	le.block = BlockPosition(block.Get())
	le.position = RecordPosition(position.Get())

	return data1
}
