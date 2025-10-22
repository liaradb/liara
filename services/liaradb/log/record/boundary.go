package record

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type Boundary struct {
	crc    CRC
	length Length
}

const BoundarySize = CrcSize + LengthSize

func NewBoundary(d []byte) Boundary {
	return Boundary{
		crc:    NewCRC(d),
		length: NewLength(uint32(len(d))),
	}
}

func (b Boundary) CRC() CRC       { return b.crc }
func (b Boundary) Length() Length { return b.length }

func (b Boundary) Size() int {
	return b.crc.Size() + b.length.Size()
}

func (b Boundary) Write(w io.Writer) error {
	return raw.WriteAll(w, b.crc, b.length)
}

func (b *Boundary) Read(r io.Reader) error {
	if err := raw.ReadAll(r, &b.crc, &b.length); err != nil {
		return err
	}

	if b.length.Value() == 0 {
		return io.EOF
	}

	return nil
}
