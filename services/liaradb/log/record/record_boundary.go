package record

import (
	"bufio"
	"io"
)

type RecordBoundary struct {
	crc    CRC
	length RecordLength
}

const RecordHeaderSize = CrcSize + RecordLengthSize

func NewRecordBoundary(d []byte) RecordBoundary {
	return RecordBoundary{
		crc:    NewCRC(d),
		length: NewRecordLength(d),
	}
}

func (rb RecordBoundary) Size() int {
	return rb.crc.Size() + rb.length.Size()
}

func (rb RecordBoundary) Write(w io.Writer) error {
	if err := rb.crc.Write(w); err != nil {
		return err
	}

	return rb.length.Write(w)
}

func (rb *RecordBoundary) Read(r io.Reader) error {
	if err := rb.crc.Read(r); err != nil {
		return err
	}

	return rb.length.Read(r)
}

func (rb *RecordBoundary) Validate(r *bufio.Reader) error {
	if err := rb.Read(r); err != nil {
		return err
	}

	if rb.length == 0 {
		return io.EOF
	}

	d, err := r.Peek(int(rb.length))
	if err != nil {
		return err
	}

	if !rb.crc.Compare(d) {
		return ErrInvalidCRC
	}

	return nil
}

// TODO: We need to rewind the length
func (rb *RecordBoundary) Skip(r io.Reader) error {
	if err := rb.Read(r); err != nil {
		return err
	}

	if rb.length == 0 {
		return io.EOF
	}

	return nil
}
