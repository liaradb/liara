package page

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

// TODO: This must be int32 to prevent overflows
type Offset int64

const OffsetSize = 8

func (o Offset) Value() int64 { return int64(o) }

func (Offset) Size() int { return OffsetSize }

func (o Offset) Write(w io.Writer) error {
	return raw.WriteInt64(w, o)
}

func (o *Offset) Read(r io.Reader) error {
	return raw.ReadInt64(r, o)
}
