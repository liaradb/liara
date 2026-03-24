package base

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

type Uint64 uint64

const Uint64Size = 8

func NewUint64(value uint64) Uint64 {
	return Uint64(value)
}

func NewInt64(value int64) Uint64 {
	return Uint64(value)
}

func (b Uint64) Value() uint64  { return uint64(b) }
func (b Uint64) Signed() int64  { return int64(b) }
func (b Uint64) Size() int      { return Uint64Size }
func (b Uint64) String() string { return fmt.Sprintf("%016x", b.Value()) }

func (b Uint64) Write(w io.Writer) error {
	return raw.WriteInt64(w, b)
}

func (b *Uint64) Read(r io.Reader) error {
	return raw.ReadInt64(r, b)
}

func (b Uint64) WriteData(data []byte) []byte {
	return scan.SetUint64(data, b.Value())
}

func (b *Uint64) ReadData(data []byte) []byte {
	v, data0 := scan.Uint64(data)
	*b = Uint64(v)
	return data0
}
