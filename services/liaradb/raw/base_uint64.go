package raw

import (
	"fmt"
	"io"
)

type BaseUint64 uint64

const BaseUint64Size = 8

func NewBaseUint64(value uint64) BaseUint64 {
	return BaseUint64(value)
}

func (b BaseUint64) Value() uint64  { return uint64(b) }
func (b BaseUint64) Size() int      { return BaseUint64Size }
func (b BaseUint64) String() string { return fmt.Sprintf("%v", b.Value()) }

func (b BaseUint64) Write(w io.Writer) error {
	return WriteInt64(w, b)
}

func (b *BaseUint64) Read(r io.Reader) error {
	return ReadInt64(r, b)
}
