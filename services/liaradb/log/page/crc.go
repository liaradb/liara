package page

import (
	"bufio"
	"encoding/binary"
	"hash/crc32"
	"io"
)

type CRC uint32

const CrcSize = 4

var table = crc32.MakeTable(crc32.Castagnoli)

func NewCRC(data []byte) CRC {
	return CRC(crc32.Checksum(data, table))
}

func (c CRC) Compare(data []byte) bool {
	return NewCRC(data) == c
}

func (c CRC) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c)
}

func (c *CRC) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, c)
}

const RecordHeaderSize = CrcSize + RecordLengthSize

func WriteCRC(crc CRC, data []byte, w *bufio.Writer) error {
	if err := crc.Write(w); err != nil {
		return err
	}

	if err := NewRecordLength(data).Write(w); err != nil {
		return err
	}

	return nil
}

func ValidateCRC(r *bufio.Reader) error {
	var c CRC
	if err := c.Read(r); err != nil {
		return err
	}

	rl := RecordLength(0)
	if err := rl.Read(r); err != nil {
		return err
	}

	if rl == 0 {
		return io.EOF
	}

	d, err := r.Peek(int(rl))
	if err != nil {
		return err
	}

	if !c.Compare(d) {
		return ErrInvalidCRC
	}

	return nil
}

// TODO: We need to rewind the length
func SkipCRC(r io.Reader) error {
	var c CRC
	if err := c.Read(r); err != nil {
		return err
	}

	rl := RecordLength(0)
	if err := rl.Read(r); err != nil {
		return err
	}

	if rl == 0 {
		return io.EOF
	}

	return nil
}
