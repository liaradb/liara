package page

import (
	"encoding/binary"
	"io"
)

type Offset uint32

const OffsetSize = 4

func (Offset) Size() int { return OffsetSize }

func (o Offset) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, o)
}

func (o *Offset) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, o)
}
