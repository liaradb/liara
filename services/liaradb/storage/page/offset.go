package page

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

// TODO: This must be int32 to prevent overflows
type Offset int32

const OffsetSize = 4

func (o Offset) Value() int { return int(o) }

func (Offset) Size() int { return OffsetSize }

func (o Offset) Write(w io.Writer) error {
	return raw.WriteInt32(w, o)
}

func (o *Offset) Read(r io.Reader) error {
	return raw.ReadInt32(r, o)
}
