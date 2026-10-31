package link

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

type BlockID int64

const BlockIDSize = 8

func (p BlockID) Value() int64   { return int64(p) }
func (BlockID) Size() int        { return BlockIDSize }
func (p BlockID) String() string { return fmt.Sprintf("%v", p.Value()) }

func (p BlockID) Offset(bufferSize int64) Offset {
	return Offset(p) * Offset(bufferSize)
}

func (p BlockID) Write(w io.Writer) error {
	return raw.WriteInt64(w, p)
}

func (p *BlockID) Read(r io.Reader) error {
	return raw.ReadInt64(r, p)
}

func (p BlockID) WriteData(data []byte) []byte {
	return scan.SetInt64(data, p.Value())
}

func (p *BlockID) ReadData(data []byte) []byte {
	block, data0 := scan.Int64(data)
	*p = BlockID(block)
	return data0
}
