package link

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

type RecordPosition int16

const RecordPositionSize = 2

func (b RecordPosition) Value() int16   { return int16(b) }
func (b RecordPosition) Size() int      { return RecordPositionSize }
func (b RecordPosition) String() string { return fmt.Sprintf("%v", b.Value()) }

func (b RecordPosition) Write(w io.Writer) error {
	return raw.WriteInt16(w, b)
}

func (b *RecordPosition) Read(r io.Reader) error {
	return raw.ReadInt16(r, b)
}

func (b RecordPosition) WriteData(data []byte) ([]byte, bool) {
	return scan.SetInt16(data, b.Value())
}

func (b *RecordPosition) ReadData(data []byte) ([]byte, bool) {
	block, data0, ok := scan.Int16(data)
	if !ok {
		return nil, false
	}

	*b = RecordPosition(block)
	return data0, true
}
