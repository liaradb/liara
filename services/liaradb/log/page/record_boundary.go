package page

import (
	"bufio"
	"io"
)

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
