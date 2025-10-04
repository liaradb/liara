package record

import (
	"encoding/binary"
	"io"
)

type Length uint32

const LengthSize = 4

func NewLength(d []byte) Length {
	return Length(len(d))
}

func (l Length) Size() int {
	return LengthSize
}

func (l Length) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, l)
}

func (l *Length) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, l)
}
