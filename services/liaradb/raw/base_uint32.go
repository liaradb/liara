package raw

import (
	"fmt"
	"io"
)

type BaseUint32 uint32

const BaseUint32Size = 4

func NewBaseUint32(value uint32) BaseUint32 {
	return BaseUint32(value)
}

func (b BaseUint32) Value() uint32  { return uint32(b) }
func (b BaseUint32) Size() int      { return BaseUint32Size }
func (b BaseUint32) String() string { return fmt.Sprintf("%v", b.Value()) }

func (b BaseUint32) Write(w io.Writer) error {
	return WriteInt32(w, b)
}

func (b *BaseUint32) Read(r io.Reader) error {
	return ReadInt32(r, b)
}
