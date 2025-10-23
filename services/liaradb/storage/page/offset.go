package page

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type Offset uint32

const OffsetSize = 4

func (Offset) Size() int { return OffsetSize }

func (o Offset) Write(w io.Writer) error {
	return raw.WriteInt32(w, o)
}

func (o *Offset) Read(r io.Reader) error {
	return raw.ReadInt32(r, o)
}
