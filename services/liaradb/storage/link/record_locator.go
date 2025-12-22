package link

import (
	"github.com/liaradb/liaradb/encoder/wrap"
)

const RecordIDSize = 8 + 1

// TODO: Test this
type RecordLocator struct {
	block    FilePosition
	position RecordPosition
}

func NewRecordLocator(block FilePosition, position RecordPosition) RecordLocator {
	return RecordLocator{
		block:    block,
		position: position,
	}
}

func (i RecordLocator) Block() FilePosition { return i.block }
func (i RecordLocator) Position() int8      { return i.position.Value() }
func (i RecordLocator) Size() int           { return RecordIDSize }

// TODO: Use simpler wrap
func (le RecordLocator) Write(data []byte) []byte {
	block, data0 := wrap.NewInt64(data)
	position, data1 := wrap.NewByte(data0)

	block.Set(le.block.Value())
	position.Set(le.position.Value())

	return data1
}

// TODO: Use simpler wrap
func (le *RecordLocator) Read(data []byte) []byte {
	block, data0 := wrap.NewInt64(data)
	position, data1 := wrap.NewByte(data0)

	le.block = FilePosition(block.Get())
	le.position = RecordPosition(position.Get())

	return data1
}
