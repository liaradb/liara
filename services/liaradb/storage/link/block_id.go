package link

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

type BlockID int64

const BlockIDSize = 8

func (i BlockID) Value() int64   { return int64(i) }
func (BlockID) Size() int        { return BlockIDSize }
func (i BlockID) String() string { return fmt.Sprintf("%v", i.Value()) }

func (i BlockID) Offset(bufferSize int64) Offset {
	return Offset(i) * Offset(bufferSize)
}

func (i BlockID) Write(w io.Writer) error {
	return raw.WriteInt64(w, i)
}

func (i *BlockID) Read(r io.Reader) error {
	return raw.ReadInt64(r, i)
}

func (i BlockID) WriteData(data []byte) []byte {
	return scan.SetInt64(data, i.Value())
}

func (i *BlockID) ReadData(data []byte) []byte {
	block, data0 := scan.Int64(data)
	*i = BlockID(block)
	return data0
}
