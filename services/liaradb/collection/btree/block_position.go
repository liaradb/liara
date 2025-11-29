package btree

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

// TODO: Test this
type BlockPosition uint64

const BlockPositionSize = 8

func (b BlockPosition) Value() int64   { return int64(b) }
func (b BlockPosition) Size() int      { return BlockPositionSize }
func (b BlockPosition) String() string { return fmt.Sprintf("%v", b.Value()) }

func (b BlockPosition) Write(w io.Writer) error {
	return raw.WriteInt64(w, b)
}

func (b *BlockPosition) Read(r io.Reader) error {
	return raw.ReadInt64(r, b)
}
