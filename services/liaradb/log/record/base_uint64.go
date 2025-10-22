package record

import (
	"encoding/binary"
	"fmt"
	"io"
)

type baseUint64 uint64

const baseUint64Size = 8

func NewBaseUint64(value uint64) baseUint64 {
	return baseUint64(value)
}

func (b baseUint64) Value() uint64  { return uint64(b) }
func (b baseUint64) Size() int      { return baseUint64Size }
func (b baseUint64) String() string { return fmt.Sprintf("%v", b.Value()) }

func (b baseUint64) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, b)
}

func (b *baseUint64) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, b)
}
