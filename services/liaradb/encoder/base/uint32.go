package base

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

type Uint32 uint32

const Uint32Size = 4

func NewUint32(value uint32) Uint32 {
	return Uint32(value)
}

func (b Uint32) Value() uint32  { return uint32(b) }
func (b Uint32) Size() int      { return Uint32Size }
func (b Uint32) String() string { return fmt.Sprintf("%08x", b.Value()) }

func (b Uint32) Write(w io.Writer) error {
	return raw.WriteInt32(w, b)
}

func (b *Uint32) Read(r io.Reader) error {
	return raw.ReadInt32(r, b)
}

func (b Uint32) WriteData(data []byte) ([]byte, bool) {
	return scan.SetUint32(data, b.Value())
}

func (b *Uint32) ReadData(data []byte) ([]byte, bool) {
	v, data0, ok := scan.Uint32(data)
	if !ok {
		return nil, false
	}

	*b = Uint32(v)
	return data0, true
}
