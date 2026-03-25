package link

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

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

func (b RecordPosition) WriteData(data []byte) ([]byte, bool) {
	return scan.SetInt8(data, b.Value())
}

func (b *RecordPosition) ReadData(data []byte) ([]byte, bool) {
	block, data0, ok := scan.Int8(data)
	if !ok {
		return nil, false
	}

	*b = RecordPosition(block)
	return data0, true
}
