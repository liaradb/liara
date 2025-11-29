package btree

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
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

func (le RecordID) Write(w io.Writer) error {
	return raw.WriteAll(w,
		le.block,
		le.position,
	)
}

func (le *RecordID) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&le.block,
		&le.position)
}
