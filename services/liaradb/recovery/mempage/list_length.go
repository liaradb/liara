package mempage

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type listLength uint32

const ListLengthSize = 4

func (listLength) Size() int     { return ListLengthSize }
func (ll listLength) Value() int { return int(ll) }

func (ll listLength) Write(w io.Writer) error {
	return raw.WriteInt32(w, ll)
}

func (ll *listLength) Read(r io.Reader) error {
	return raw.ReadInt32(r, ll)
}
