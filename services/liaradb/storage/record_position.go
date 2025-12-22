package storage

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

// TODO: Test this
type RecordPosition int8

const RecordPositionSize = 1

func (b RecordPosition) Value() int8    { return int8(b) }
func (b RecordPosition) Size() int      { return RecordPositionSize }
func (b RecordPosition) String() string { return fmt.Sprintf("%v", b.Value()) }

func (b RecordPosition) Write(w io.Writer) error {
	return raw.WriteInt8(w, b)
}

func (b *RecordPosition) Read(r io.Reader) error {
	return raw.ReadInt8(r, b)
}
