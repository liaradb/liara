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

func (le RecordLocator) Write(data []byte) ([]byte, bool) {
	data0, ok := scan.SetInt64(data, le.block.Value())
	if !ok {
		return nil, false
	}

	data1 := scan.SetInt8(data0, le.position.Value())

	return data1, true
}

func (le *RecordLocator) Read(data []byte) ([]byte, bool) {
	block, data0, ok := scan.Int64(data)
	if !ok {
		return nil, false
	}

	position, data1 := scan.Int8(data0)

	le.block = FilePosition(block)
	le.position = RecordPosition(position)

	return data1, true
}
