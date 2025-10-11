package record

import (
	"encoding/binary"
	"io"
)

type ListLength uint32

const ListLengthSize = 4

func (ListLength) Size() int { return ListLengthSize }

func (ll ListLength) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, ll)
}

func (ll *ListLength) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, ll)
}
