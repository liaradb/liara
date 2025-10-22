package record

import (
	"io"
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
	if err := b.crc.Write(w); err != nil {
		return err
	}

	return b.length.Write(w)
}

func (b *Boundary) Read(r io.Reader) error {
	if err := b.crc.Read(r); err != nil {
		return err
	}

	if err := b.length.Read(r); err != nil {
		return err
	}

	if b.length.Value() == 0 {
		return io.EOF
	}

	return nil
}
