package record

import (
	"encoding/binary"
	"io"
)

type baseUint32 uint32

const baseUint32Size = 4

func NewBaseUint32(size uint32) baseUint32 {
	return baseUint32(size)
}

func (b baseUint32) Value() uint32 { return uint32(b) }
func (b baseUint32) Size() int     { return baseUint32Size }

func (b baseUint32) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, b)
}

func (b *baseUint32) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, b)
}
