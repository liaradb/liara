package base

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

type BaseUint64 uint64

const BaseUint64Size = 8

func NewBaseUint64(value uint64) BaseUint64 {
	return BaseUint64(value)
}

func (b BaseUint64) Value() uint64  { return uint64(b) }
func (b BaseUint64) Size() int      { return BaseUint64Size }
func (b BaseUint64) String() string { return fmt.Sprintf("%016x", b.Value()) }

func (b BaseUint64) Write(w io.Writer) error {
	return raw.WriteInt64(w, b)
}

func (b *BaseUint64) Read(r io.Reader) error {
	return raw.ReadInt64(r, b)
}

func (b BaseUint64) WriteData(data []byte) []byte {
	return scan.SetUint64(data, b.Value())
}

func (b *BaseUint64) ReadData(data []byte) []byte {
	v, data0 := scan.Uint64(data)
	*b = BaseUint64(v)
	return data0
}
