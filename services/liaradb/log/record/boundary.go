package record

import (
	"bufio"
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
		length: NewLength(d),
	}
}

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

	return b.length.Read(r)
}

// TODO: This reads past the end of the file
func (b *Boundary) Validate(r *bufio.Reader) error {
	if err := b.Read(r); err != nil {
		return err
	}

	if b.length == 0 {
		return io.EOF
	}

	d, err := r.Peek(int(b.length))
	if err != nil {
		return err
	}

	if !b.crc.Compare(d) {
		return ErrInvalidCRC
	}

	return nil
}

// TODO: We need to rewind the length
func (b *Boundary) Skip(r io.Reader) error {
	if err := b.Read(r); err != nil {
		return err
	}

	if b.length == 0 {
		return io.EOF
	}

	return nil
}
