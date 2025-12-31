package mempage

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type ListLength uint32

const ListLengthSize = 4

func (ListLength) Size() int     { return ListLengthSize }
func (ll ListLength) Value() int { return int(ll) }

func (ll ListLength) Write(w io.Writer) error {
	return raw.WriteInt32(w, ll)
}

func (ll *ListLength) Read(r io.Reader) error {
	return raw.ReadInt32(r, ll)
}
