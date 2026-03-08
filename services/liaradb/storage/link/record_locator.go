package link

import (
	"github.com/liaradb/liaradb/encoder/scan"
)

const RecordLocatorSize = 8 + 1

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
func (i RecordLocator) Size() int           { return RecordLocatorSize }

func (le RecordLocator) Write(data []byte) []byte {
	data0 := scan.SetInt64(data, le.block.Value())
	data1 := scan.SetInt8(data0, le.position.Value())

	return data1
}

func (le *RecordLocator) Read(data []byte) []byte {
	block, data0 := scan.Int64(data)
	position, data1 := scan.Int8(data0)

	le.block = FilePosition(block)
	le.position = RecordPosition(position)

	return data1
}
